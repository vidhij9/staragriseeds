package service

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Services struct {
	Farmer *FarmerService
	CCE    *CCEService
	Ticket *TicketService
	Shoot  *ShootService
}

func NewServices(dbClient *dynamodb.Client) *Services {
	return &Services{
		Farmer: NewFarmerService(dbClient),
		CCE:    NewCCEService(dbClient),
		Ticket: NewTicketService(dbClient),
	}
}
