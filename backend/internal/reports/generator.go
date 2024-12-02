// In internal/reports/generator.go

package reports

import (
	"backend/internal/models"
	"backend/internal/service"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type ReportGenerator struct {
	shootService *service.ShootService
	cceService   *service.CCEService
	dbClient     *dynamodb.Client
}

func NewReportGenerator(shootService *service.ShootService, cceService *service.CCEService, dbClient *dynamodb.Client) *ReportGenerator {
	return &ReportGenerator{
		shootService: shootService,
		cceService:   cceService,
		dbClient:     dbClient,
	}
}

func (rg *ReportGenerator) GenerateReport(ctx context.Context, reportType string, startDate, endDate time.Time) (*models.Report, error) {
	shoots, err := rg.shootService.GetShootsWithDateFilter(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	cces, _, err := rg.cceService.ListCCEs(ctx, "")
	if err != nil {
		return nil, err
	}

	report := &models.Report{
		ID:         uuid.New().String(),
		ReportType: reportType,
		StartDate:  startDate,
		EndDate:    endDate,
		CreatedAt:  time.Now(),
		CCEReports: make(map[string]*models.CCEReport),
	}

	for _, cce := range cces {
		report.CCEReports[cce.ID] = &models.CCEReport{
			CCEID: cce.ID,
			Name:  cce.Name,
		}
	}

	for _, shoot := range shoots {
		cceReport := report.CCEReports[shoot.CCEID]
		if cceReport == nil {
			continue
		}

		if shoot.Status == "missed" {
			cceReport.MissedCalls++
		} else {
			cceReport.CompletedShoots++
			if shoot.Type == "call" {
				cceReport.AttendedCalls++
				cceReport.TotalTalkTime += shoot.Duration
			}
		}
	}

	for _, cceReport := range report.CCEReports {
		if cceReport.AttendedCalls > 0 {
			cceReport.AvgTalkTime = float64(cceReport.TotalTalkTime) / float64(cceReport.AttendedCalls)
		}
	}

	return report, nil
}

func (rg *ReportGenerator) SaveReport(ctx context.Context, report *models.Report) error {
	// Assuming you're using DynamoDB to store reports
	item, err := attributevalue.MarshalMap(report)
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	_, err = rg.dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("Reports"),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to save report: %w", err)
	}

	return nil
}
