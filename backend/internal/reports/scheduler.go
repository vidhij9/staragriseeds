package reports

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron            *cron.Cron
	reportGenerator *ReportGenerator
	mailer          *Mailer
}

func NewScheduler(reportGenerator *ReportGenerator, mailer *Mailer) *Scheduler {
	return &Scheduler{
		cron:            cron.New(),
		reportGenerator: reportGenerator,
		mailer:          mailer,
	}
}

func (s *Scheduler) Start() {
	// Schedule daily report at 1:00 AM UTC
	s.cron.AddFunc("0 1 * * *", func() {
		s.generateAndSendReport("daily")
	})

	// Schedule weekly report at 2:00 AM UTC on Mondays
	s.cron.AddFunc("0 2 * * 1", func() {
		s.generateAndSendReport("weekly")
	})

	// Schedule monthly report at 3:00 AM UTC on the 1st of each month
	s.cron.AddFunc("0 3 1 * *", func() {
		s.generateAndSendReport("monthly")
	})

	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}

func (s *Scheduler) generateAndSendReport(reportType string) {
	ctx := context.Background()

	var startDate, endDate time.Time
	now := time.Now().UTC()

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
		log.Print("Invalid report type", "type", reportType)
		return
	}

	report, err := s.reportGenerator.GenerateReport(ctx, reportType, startDate, endDate)
	if err != nil {
		log.Print("Failed to generate report", "type", reportType, "error", err)
		return
	}

	if err := s.reportGenerator.SaveReport(ctx, report); err != nil {
		log.Print("Failed to save report", "type", reportType, "error", err)
		return
	}

	if err := s.mailer.SendReport(ctx, report); err != nil {
		log.Print("Failed to send report", "type", reportType, "error", err)
		return
	}

	log.Print("Report generated and sent successfully", "type", reportType)
}
