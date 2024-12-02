package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Migration struct {
	Version     int
	Description string
	Up          func(context.Context, *dynamodb.Client) error
	Down        func(context.Context, *dynamodb.Client) error
}

var migrations = []Migration{
	{
		Version:     1,
		Description: "Create initial tables",
		Up: func(ctx context.Context, client *dynamodb.Client) error {
			if err := createTable(ctx, client, "Farmers"); err != nil {
				return err
			}
			if err := createTable(ctx, client, "CCEs"); err != nil {
				return err
			}
			if err := createTable(ctx, client, "Tickets"); err != nil {
				return err
			}
			return nil
		},
		Down: func(ctx context.Context, client *dynamodb.Client) error {
			if err := deleteTable(ctx, client, "Farmers"); err != nil {
				return err
			}
			if err := deleteTable(ctx, client, "CCEs"); err != nil {
				return err
			}
			if err := deleteTable(ctx, client, "Tickets"); err != nil {
				return err
			}
			return nil
		},
	},
	// Add more migrations here as your schema evolves
}

func RunMigrations(client *dynamodb.Client) error {
	ctx := context.Background()

	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(ctx, client); err != nil {
		return err
	}

	// Get current migration version
	currentVersion, err := getCurrentMigrationVersion(ctx, client)
	if err != nil {
		return err
	}

	// Run pending migrations
	for _, migration := range migrations {
		if migration.Version > currentVersion {
			log.Printf("Running migration %d: %s", migration.Version, migration.Description)
			if err := migration.Up(ctx, client); err != nil {
				return fmt.Errorf("failed to run migration %d: %w", migration.Version, err)
			}
			if err := updateMigrationVersion(ctx, client, migration.Version); err != nil {
				return fmt.Errorf("failed to update migration version: %w", err)
			}
		}
	}

	return nil
}

func createMigrationsTable(ctx context.Context, client *dynamodb.Client) error {
	_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String("Migrations"),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("Key"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("Key"),
				KeyType:       types.KeyTypeHash,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	})
	if err != nil {
		var resourceInUseErr *types.ResourceInUseException
		if ok := errors.As(err, &resourceInUseErr); ok {
			// Table already exists, ignore the error
			return nil
		}
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	return nil
}

func getCurrentMigrationVersion(ctx context.Context, client *dynamodb.Client) (int, error) {
	result, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("Migrations"),
		Key: map[string]types.AttributeValue{
			"Key": &types.AttributeValueMemberS{Value: "Version"},
		},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get current migration version: %w", err)
	}

	if result.Item == nil {
		return 0, nil
	}

	var version struct {
		Version int `dynamodbav:"Version"`
	}
	if err := attributevalue.UnmarshalMap(result.Item, &version); err != nil {
		return 0, fmt.Errorf("failed to unmarshal migration version: %w", err)
	}

	return version.Version, nil
}

func updateMigrationVersion(ctx context.Context, client *dynamodb.Client, version int) error {
	item, err := attributevalue.MarshalMap(struct {
		Key     string `dynamodbav:"Key"`
		Version int    `dynamodbav:"Version"`
	}{
		Key:     "Version",
		Version: version,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal migration version: %w", err)
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("Migrations"),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to update migration version: %w", err)
	}

	return nil
}

func createTable(ctx context.Context, client *dynamodb.Client, tableName string) error {
	_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       types.KeyTypeHash,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create %s table: %w", tableName, err)
	}
	log.Printf("Table %s created successfully", tableName)
	return nil
}

func deleteTable(ctx context.Context, client *dynamodb.Client, tableName string) error {
	_, err := client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete table %s: %w", tableName, err)
	}
	log.Printf("Table %s deleted successfully", tableName)
	return nil
}
