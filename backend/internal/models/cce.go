package models

type CCE struct {
	ID      string   `json:"id" dynamodbav:"id"`
	Name    string   `json:"name" dynamodbav:"name"`
	Farmers []Farmer `json:"farmers" dynamodbav:"farmers"`
	Tickets []Ticket `json:"tickets" dynamodbav:"tickets"`
	AvgTime float64  `json:"avgTime" dynamodbav:"avgTime"`
}
