package handlers

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DynamoDB client
var dbClient *dynamodb.Client

// InitDynamoDB initializes the DynamoDB client
func InitDynamoDB(cfg aws.Config) {
	dbClient = dynamodb.NewFromConfig(cfg)
}
