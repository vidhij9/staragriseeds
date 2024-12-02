package service

import (
	"context"
	"sort"
	"time"

	"backend/internal/models"
	"backend/pkg/errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const TicketTableName = "Tickets"

type TicketService struct {
	dbClient *dynamodb.Client
}

func NewTicketService(dbClient *dynamodb.Client) *TicketService {
	return &TicketService{
		dbClient: dbClient,
	}
}

func (s *TicketService) CreateTicket(ctx context.Context, ticket *models.Ticket) error {
	item, err := attributevalue.MarshalMap(ticket)
	if err != nil {
		return errors.ErrInternal
	}

	_, err = s.dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(TicketTableName),
		Item:      item,
	})
	if err != nil {
		return errors.ErrInternal
	}

	return nil
}

func (s *TicketService) GetTicket(ctx context.Context, id string) (*models.Ticket, error) {
	result, err := s.dbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(TicketTableName),
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

	var ticket models.Ticket
	err = attributevalue.UnmarshalMap(result.Item, &ticket)
	if err != nil {
		return nil, errors.ErrInternal
	}

	return &ticket, nil
}

func (s *TicketService) GetTicketsByFarmerContact(ctx context.Context, farmer *models.Farmer) ([]models.Ticket, error) {
	// farmer, err := s.farmerService.GetFarmerByContact(ctx, contact)
	// if err != nil {
	// 	return nil, err
	// }

	result, err := s.dbClient.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String("Tickets"),
		IndexName:              aws.String("FarmerIDIndex"),
		KeyConditionExpression: aws.String("FarmerID = :farmerID"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":farmerID": &types.AttributeValueMemberS{Value: farmer.ID},
		},
	})
	if err != nil {
		return nil, err
	}

	var tickets []models.Ticket
	err = attributevalue.UnmarshalListOfMaps(result.Items, &tickets)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (s *TicketService) GetTicketsByCCE(ctx context.Context, cceID string) ([]models.Ticket, error) {
	result, err := s.dbClient.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String("Tickets"),
		IndexName:              aws.String("CCEIDIndex"),
		KeyConditionExpression: aws.String("CCEID = :cceID"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":cceID": &types.AttributeValueMemberS{Value: cceID},
		},
	})
	if err != nil {
		return nil, err
	}

	var tickets []models.Ticket
	err = attributevalue.UnmarshalListOfMaps(result.Items, &tickets)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (s *TicketService) GetTicketsByCCEWithDateFilter(ctx context.Context, cceID string, startDate, endDate time.Time) ([]models.Ticket, error) {
	result, err := s.dbClient.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String("Tickets"),
		IndexName:              aws.String("CCEIDCreatedAtIndex"),
		KeyConditionExpression: aws.String("CCEID = :cceID AND CreatedAt BETWEEN :startDate AND :endDate"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":cceID":     &types.AttributeValueMemberS{Value: cceID},
			":startDate": &types.AttributeValueMemberS{Value: startDate.Format(time.RFC3339)},
			":endDate":   &types.AttributeValueMemberS{Value: endDate.Format(time.RFC3339)},
		},
	})
	if err != nil {
		return nil, err
	}

	var tickets []models.Ticket
	err = attributevalue.UnmarshalListOfMaps(result.Items, &tickets)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (s *TicketService) GetTicketsByCCEAndStatus(ctx context.Context, cceID, status string) ([]models.Ticket, error) {
	result, err := s.dbClient.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String("Tickets"),
		IndexName:              aws.String("CCEIDStatusIndex"),
		KeyConditionExpression: aws.String("CCEID = :cceID AND #status = :status"),
		ExpressionAttributeNames: map[string]string{
			"#status": "Status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":cceID":  &types.AttributeValueMemberS{Value: cceID},
			":status": &types.AttributeValueMemberS{Value: status},
		},
	})
	if err != nil {
		return nil, err
	}

	var tickets []models.Ticket
	err = attributevalue.UnmarshalListOfMaps(result.Items, &tickets)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (s *TicketService) GetTicketsWithStatusAndSort(ctx context.Context, status, sortBy, sortOrder string) ([]models.Ticket, error) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String("Tickets"),
		FilterExpression: aws.String("#status = :status"),
		ExpressionAttributeNames: map[string]string{
			"#status": "Status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status": &types.AttributeValueMemberS{Value: status},
		},
	}

	result, err := s.dbClient.Scan(ctx, input)
	if err != nil {
		return nil, err
	}

	var tickets []models.Ticket
	err = attributevalue.UnmarshalListOfMaps(result.Items, &tickets)
	if err != nil {
		return nil, err
	}

	// Sort the tickets based on sortBy and sortOrder
	if sortBy != "" {
		sort.Slice(tickets, func(i, j int) bool {
			switch sortBy {
			case "createdAt":
				if sortOrder == "desc" {
					return tickets[i].CreatedAt.After(tickets[j].CreatedAt)
				}
				return tickets[i].CreatedAt.Before(tickets[j].CreatedAt)
			case "updatedAt":
				if sortOrder == "desc" {
					return tickets[i].UpdatedAt.After(tickets[j].UpdatedAt)
				}
				return tickets[i].UpdatedAt.Before(tickets[j].UpdatedAt)
			default:
				return false
			}
		})
	}

	return tickets, nil
}

func (s *TicketService) UpdateTicket(ctx context.Context, ticket *models.Ticket) error {
	item, err := attributevalue.MarshalMap(ticket)
	if err != nil {
		return errors.ErrInternal
	}

	_, err = s.dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(TicketTableName),
		Item:      item,
	})
	if err != nil {
		return errors.ErrInternal
	}

	return nil
}

func (s *TicketService) DeleteTicket(ctx context.Context, id string) error {
	_, err := s.dbClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(TicketTableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return errors.ErrInternal
	}

	return nil
}

func (s *TicketService) ListTickets(ctx context.Context, limit int32, nextToken string) ([]models.Ticket, string, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(TicketTableName),
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

	var tickets []models.Ticket
	err = attributevalue.UnmarshalListOfMaps(result.Items, &tickets)
	if err != nil {
		return nil, "", errors.ErrInternal
	}

	var newNextToken string
	if result.LastEvaluatedKey != nil {
		newNextToken = result.LastEvaluatedKey["ID"].(*types.AttributeValueMemberS).Value
	}

	return tickets, newNextToken, nil
}
