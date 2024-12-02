package models

import "time"

type Shoot struct {
	ID        string    `json:"id" dynamodbav:"ID"`
	FarmerID  string    `json:"farmerId" dynamodbav:"FarmerID"`
	CCEID     string    `json:"cceId" dynamodbav:"CCEID"`
	Type      string    `json:"type" dynamodbav:"Type"`     // "whatsapp" or "call"
	Status    string    `json:"status" dynamodbav:"Status"` // "completed" or "missed"
	Timestamp time.Time `json:"timestamp" dynamodbav:"Timestamp"`
	Duration  int       `json:"duration" dynamodbav:"Duration"` // in seconds
}

type TypeOfShoot struct {
	WhatsAppSent bool `json:"whatsapp_sent" dynamodbav:"whatsapp_sent"`
	CallBack     bool `json:"call_back" dynamodbav:"call_back"`
}
