package reports

import (
	"backend/internal/models"
	"context"
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	dialer *gomail.Dialer
}

func NewMailer(host string, port int, username, password string) *Mailer {
	return &Mailer{
		dialer: gomail.NewDialer(host, port, username, password),
	}
}

func (m *Mailer) SendReport(ctx context.Context, report *models.Report) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", "reports@yourcompany.com")
	msg.SetHeader("To", "recipient@example.com")
	msg.SetHeader("Subject", fmt.Sprintf("%s Report: %s to %s",
		report.ReportType,
		report.StartDate.Format("2006-01-02"),
		report.EndDate.Format("2006-01-02")))

	// Generate HTML content for the report
	content := generateReportContent(report)
	msg.SetBody("text/html", content)

	if err := m.dialer.DialAndSend(msg); err != nil {
		log.Print("Failed to send email", "error", err)
		return err
	}

	log.Print("Report sent successfully", "report_id", report.ID)
	return nil
}

func generateReportContent(report *models.Report) string {
	// Implement this function to generate HTML content for the report
	// You can use a template engine like html/template for more complex reports
	// For now, we'll use a simple string
	return fmt.Sprintf(`
        <h1>%s Report</h1>
        <p>Period: %s to %s</p>
        <p>Total CCEs: %d</p>
        <!-- Add more report details here -->
    `, report.ReportType, report.StartDate.Format("2006-01-02"), report.EndDate.Format("2006-01-02"), len(report.CCEReports))
}
