package service

import (
	"backend/internal/models"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const ShootTableName = "Shoots"

type ShootService struct {
	dbClient *dynamodb.Client
}

func NewShootService(dbClient *dynamodb.Client) *ShootService {
	return &ShootService{
		dbClient: dbClient,
	}
}

func (s *ShootService) GetAllShoots(ctx context.Context, shootType string) ([]models.Shoot, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("Shoots"),
	}

	if shootType != "" {
		input.FilterExpression = aws.String("#type = :type")
		input.ExpressionAttributeNames = map[string]string{
			"#type": "Type",
		}
		input.ExpressionAttributeValues = map[string]types.AttributeValue{
			":type": &types.AttributeValueMemberS{Value: shootType},
		}
	}

	result, err := s.dbClient.Scan(ctx, input)
	if err != nil {
		return nil, err
	}

	var shoots []models.Shoot
	err = attributevalue.UnmarshalListOfMaps(result.Items, &shoots)
	if err != nil {
		return nil, err
	}

	return shoots, nil
}

func (s *ShootService) GetShootsWithDateFilter(ctx context.Context, startDate, endDate time.Time) ([]models.Shoot, error) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String("Shoots"),
		FilterExpression: aws.String("Timestamp BETWEEN :startDate AND :endDate"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":startDate": &types.AttributeValueMemberS{Value: startDate.Format(time.RFC3339)},
			":endDate":   &types.AttributeValueMemberS{Value: endDate.Format(time.RFC3339)},
		},
	}

	result, err := s.dbClient.Scan(ctx, input)
	if err != nil {
		return nil, err
	}

	var shoots []models.Shoot
	err = attributevalue.UnmarshalListOfMaps(result.Items, &shoots)
	if err != nil {
		return nil, err
	}

	return shoots, nil
}

func (s *ShootService) GetMissedShoots(ctx context.Context) ([]models.Shoot, error) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String("Shoots"),
		FilterExpression: aws.String("#status = :status AND #type = :type"),
		ExpressionAttributeNames: map[string]string{
			"#status": "Status",
			"#type":   "Type",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status": &types.AttributeValueMemberS{Value: "missed"},
			":type":   &types.AttributeValueMemberS{Value: "call"},
		},
	}

	result, err := s.dbClient.Scan(ctx, input)
	if err != nil {
		return nil, err
	}

	var shoots []models.Shoot
	err = attributevalue.UnmarshalListOfMaps(result.Items, &shoots)
	if err != nil {
		return nil, err
	}

	return shoots, nil
}
