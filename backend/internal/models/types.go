package models

import (
	"time"
)

type Report struct {
	ID         string                `json:"id"`
	ReportType string                `json:"reportType"`
	StartDate  time.Time             `json:"startDate"`
	EndDate    time.Time             `json:"endDate"`
	CreatedAt  time.Time             `json:"createdAt"`
	CCEReports map[string]*CCEReport `json:"cceReports"`
}

type CCEReport struct {
	CCEID           string  `json:"cceId"`
	Name            string  `json:"name"`
	MissedCalls     int     `json:"missedCalls"`
	CompletedShoots int     `json:"completedShoots"`
	AttendedCalls   int     `json:"attendedCalls"`
	TotalTalkTime   int     `json:"totalTalkTime"`
	AvgTalkTime     float64 `json:"avgTalkTime"`
}
