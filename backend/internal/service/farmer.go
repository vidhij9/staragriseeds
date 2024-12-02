package service

import (
	"context"
	"fmt"

	"backend/internal/models"
	"backend/pkg/errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const FarmerTableName = "Farmers"

type FarmerService struct {
	dbClient *dynamodb.Client
}

func NewFarmerService(dbClient *dynamodb.Client) *FarmerService {
	return &FarmerService{
		dbClient: dbClient,
	}
}

func (s *FarmerService) CreateFarmer(ctx context.Context, farmer *models.Farmer) error {
	item, err := attributevalue.MarshalMap(farmer)
	if err != nil {
		return errors.ErrInternal
	}

	_, err = s.dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(FarmerTableName),
		Item:      item,
	})
	if err != nil {
		return errors.ErrInternal
	}

	return nil
}

func (s *FarmerService) GetFarmer(ctx context.Context, id string) (*models.Farmer, error) {
	result, err := s.dbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(FarmerTableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, errors.ErrInternal
	}
	if result.Item == nil {
		return nil, errors.ErrNotFound
	}

	var farmer models.Farmer
	err = attributevalue.UnmarshalMap(result.Item, &farmer)
	if err != nil {
		return nil, errors.ErrInternal
	}

	return &farmer, nil
}

func (s *FarmerService) GetFarmerByContact(ctx context.Context, contact string) (*models.Farmer, error) {
	result, err := s.dbClient.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String("Farmers"),
		IndexName:              aws.String("ContactIndex"),
		KeyConditionExpression: aws.String("Contact = :contact"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":contact": &types.AttributeValueMemberS{Value: contact},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, errors.New("farmer not found")
	}

	var farmer models.Farmer
	err = attributevalue.UnmarshalMap(result.Items[0], &farmer)
	if err != nil {
		return nil, err
	}

	return &farmer, nil
}

func (s *FarmerService) ListFarmersWithFilters(ctx context.Context, filters map[string]string) ([]models.Farmer, error) {
	var filterExpression string
	expressionAttributeValues := make(map[string]types.AttributeValue)
	expressionAttributeNames := make(map[string]string)

	for key, value := range filters {
		if value != "" {
			if filterExpression != "" {
				filterExpression += " AND "
			}
			filterExpression += fmt.Sprintf("#%s = :%s", key, key)
			expressionAttributeValues[":"+key] = &types.AttributeValueMemberS{Value: value}
			expressionAttributeNames["#"+key] = key
		}
	}

	input := &dynamodb.ScanInput{
		TableName:                 aws.String("Farmers"),
		FilterExpression:          aws.String(filterExpression),
		ExpressionAttributeValues: expressionAttributeValues,
		ExpressionAttributeNames:  expressionAttributeNames,
	}

	result, err := s.dbClient.Scan(ctx, input)
	if err != nil {
		return nil, err
	}

	var farmers []models.Farmer
	err = attributevalue.UnmarshalListOfMaps(result.Items, &farmers)
	if err != nil {
		return nil, err
	}

	return farmers, nil
}

func (s *FarmerService) UpdateFarmer(ctx context.Context, farmer *models.Farmer) error {
	item, err := attributevalue.MarshalMap(farmer)
	if err != nil {
		return errors.ErrInternal
	}

	_, err = s.dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(FarmerTableName),
		Item:      item,
	})
	if err != nil {
		return errors.ErrInternal
	}

	return nil
}

func (s *FarmerService) DeleteFarmer(ctx context.Context, id string) error {
	_, err := s.dbClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(FarmerTableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return errors.ErrInternal
	}

	return nil
}

func (s *FarmerService) ListFarmers(ctx context.Context, limit int32, nextToken string) ([]models.Farmer, string, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(FarmerTableName),
		Limit:     aws.Int32(limit),
	}

	if nextToken != "" {
		input.ExclusiveStartKey = map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: nextToken},
		}
	}

	result, err := s.dbClient.Scan(ctx, input)
	if err != nil {
		return nil, "", errors.ErrInternal
	}

	var farmers []models.Farmer
	err = attributevalue.UnmarshalListOfMaps(result.Items, &farmers)
	if err != nil {
		return nil, "", errors.ErrInternal
	}

	var newNextToken string
	if result.LastEvaluatedKey != nil {
		newNextToken = result.LastEvaluatedKey["ID"].(*types.AttributeValueMemberS).Value
	}

	return farmers, newNextToken, nil
}
