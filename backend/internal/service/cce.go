package service

import (
	"context"

	"backend/internal/models"
	"backend/pkg/errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const CCETableName = "CCEs"

type CCEService struct {
	dbClient *dynamodb.Client
}

func NewCCEService(dbClient *dynamodb.Client) *CCEService {
	return &CCEService{
		dbClient: dbClient,
	}
}

func (s *CCEService) CreateCCE(ctx context.Context, cce *models.CCE) error {
	item, err := attributevalue.MarshalMap(cce)
	if err != nil {
		return errors.ErrInternal
	}

	_, err = s.dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(CCETableName),
		Item:      item,
	})
	if err != nil {
		return errors.ErrInternal
	}

	return nil
}

func (s *CCEService) GetCCE(ctx context.Context, id string) (*models.CCE, error) {
	result, err := s.dbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(CCETableName),
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

	var cce models.CCE
	err = attributevalue.UnmarshalMap(result.Item, &cce)
	if err != nil {
		return nil, errors.ErrInternal
	}

	return &cce, nil
}

func (s *CCEService) UpdateCCE(ctx context.Context, cce *models.CCE) error {
	item, err := attributevalue.MarshalMap(cce)
	if err != nil {
		return errors.ErrInternal
	}

	_, err = s.dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(CCETableName),
		Item:      item,
	})
	if err != nil {
		return errors.ErrInternal
	}

	return nil
}

func (s *CCEService) DeleteCCE(ctx context.Context, id string) error {
	_, err := s.dbClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(CCETableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return errors.ErrInternal
	}

	return nil
}

func (s *CCEService) ListCCEs(ctx context.Context, nextToken string) ([]models.CCE, string, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(CCETableName),
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

	var cces []models.CCE
	err = attributevalue.UnmarshalListOfMaps(result.Items, &cces)
	if err != nil {
		return nil, "", errors.ErrInternal
	}

	var newNextToken string
	if result.LastEvaluatedKey != nil {
		newNextToken = result.LastEvaluatedKey["ID"].(*types.AttributeValueMemberS).Value
	}

	return cces, newNextToken, nil
}
