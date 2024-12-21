package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/api"
	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/reports"
	"backend/internal/service"

	"github.com/robfig/cron/v3"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {

	log.Println("0")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("1")

	// Initialize DynamoDB client
	dbClient, err := db.NewDynamoDBClient(context.Background(), cfg.AWS.Region)
	if err != nil {
		log.Fatalf("Failed to create DynamoDB client: %v", err)
	}
	log.Println("2")

	// Run the migrations
	if err := db.RunMigrations(dbClient); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")

	// Initialize services
	services := initializeServices(dbClient)

	// Set up router
	router := api.SetupRouter(services)

	// Create server
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: router,
	}

	log.Println("3")

	// Start server
	go func() {
		log.Printf("Starting server on %s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Println("4")

	// Initialize report generator
	reportGenerator := reports.NewReportGenerator(services.Shoot, services.CCE, dbClient)

	// Initialize mailer
	mailer := reports.NewMailer(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password)

	// Set up cron job
	c := cron.New()

	log.Println("5")

	// Daily report at 00:01
	_, err = c.AddFunc("1 0 * * *", func() {
		generateAndSaveReport(reportGenerator, mailer, "daily")
	})
	if err != nil {
		log.Printf("Failed to set up daily cron job: %v", err)
	}

	log.Println("6")

	// Weekly report on Monday at 00:05
	_, err = c.AddFunc("5 0 * * 1", func() {
		generateAndSaveReport(reportGenerator, mailer, "weekly")
	})
	if err != nil {
		log.Printf("Failed to set up weekly cron job: %v", err)
	}

	log.Println("7")

	// Monthly report on the 1st of each month at 00:10
	_, err = c.AddFunc("10 0 1 * *", func() {
		generateAndSaveReport(reportGenerator, mailer, "monthly")
	})
	if err != nil {
		log.Printf("Failed to set up monthly cron job: %v", err)
	}

	log.Println("8")

	c.Start()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Fatalf("Shutting down server...")
	log.Println("9")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("10")

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline
	if err := srv.Shutdown(ctx); err != nil {
		c.Stop()

		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("11")

	// Stop cron scheduler (stops accepting new jobs)
	c.Stop()

	log.Printf("Server exiting: %v", err)
}

func initializeServices(dbClient *dynamodb.Client) *service.Services {
	return &service.Services{
		Farmer: service.NewFarmerService(dbClient),
		CCE:    service.NewCCEService(dbClient),
		Ticket: service.NewTicketService(dbClient),
	}
}

func generateAndSaveReport(rg *reports.ReportGenerator, mailer *reports.Mailer, reportType string) {
	ctx := context.Background()

	var startDate, endDate time.Time
	now := time.Now().UTC()

	log.Println("Inside setuprouter")

	switch reportType {
	case "daily":
		endDate = now.Truncate(24 * time.Hour)
		startDate = endDate.Add(-24 * time.Hour)
	case "weekly":
		endDate = now.Truncate(24 * time.Hour)
		startDate = endDate.Add(-7 * 24 * time.Hour)
	case "monthly":
		endDate = now.Truncate(24 * time.Hour)
		startDate = endDate.AddDate(0, -1, 0)
	default:
		log.Fatalf("Invalid report type: %s", reportType)
		return
	}

	// Generate the report
	report, err := rg.GenerateReport(ctx, reportType, startDate, endDate)
	if err != nil {
		log.Fatalf("Failed to generate %s report: %v", reportType, err)
		return
	}

	// Save the report
	err = rg.SaveReport(ctx, report)
	if err != nil {
		log.Fatalf("Failed to save %s report: %v", reportType, err)
		return
	}

	// Send the report
	err = mailer.SendReport(ctx, report)
	if err != nil {
		log.Fatalf("Failed to send %s report: %v", reportType, err)
		return
	}

	log.Printf("Generated, saved, and sent %s report from %s to %s", reportType, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
}
