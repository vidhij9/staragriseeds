package models

import "time"

type Ticket struct {
	ID          string    `json:"id" dynamodbav:"ID"`
	FarmerID    string    `json:"farmerId" dynamodbav:"FarmerID"`
	CCEID       string    `json:"cceId" dynamodbav:"CCEID"`
	Description string    `json:"description" dynamodbav:"Description"`
	Status      string    `json:"status" dynamodbav:"Status"`
	CreatedAt   time.Time `json:"createdAt" dynamodbav:"CreatedAt"`
	UpdatedAt   time.Time `json:"updatedAt" dynamodbav:"UpdatedAt"`
}
