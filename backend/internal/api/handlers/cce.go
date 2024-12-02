package handlers

import (
	"backend/internal/models"
	"backend/internal/service"
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gorilla/mux"
)

const CceTableName = "CCEs"

type CCEHandler struct {
	cceService *service.CCEService
}

func NewCCEHandler(cceService *service.CCEService) *CCEHandler {
	return &CCEHandler{cceService: cceService}
}

// GetCCE - Retrieve CCE by ID
func (h *CCEHandler) GetCCE(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cceID := vars["id"]

	result, err := dbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(CceTableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: cceID},
		},
	})
	if err != nil || result.Item == nil {
		http.Error(w, "CCE not found", http.StatusNotFound)
		return
	}

	var cce models.CCE
	err = attributevalue.UnmarshalMap(result.Item, &cce)
	if err != nil {
		http.Error(w, "Failed to unmarshal CCE data", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(cce)
}

// GetCCEs - Retrieve all CCEs
func (h *CCEHandler) GetCCEs(w http.ResponseWriter, r *http.Request) {
	result, err := dbClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(CceTableName),
	})
	if err != nil {
		http.Error(w, "Failed to scan CCEs", http.StatusInternalServerError)
		return
	}

	var cces []models.CCE
	for _, item := range result.Items {
		var cce models.CCE
		err = attributevalue.UnmarshalMap(item, &cce)
		if err != nil {
			http.Error(w, "Failed to unmarshal CCE data", http.StatusInternalServerError)
			return
		}
		cces = append(cces, cce)
	}

	json.NewEncoder(w).Encode(cces)
}

// CreateCCE - Add new CCE
func (h *CCEHandler) CreateCCE(w http.ResponseWriter, r *http.Request) {
	var cce models.CCE
	err := json.NewDecoder(r.Body).Decode(&cce)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err = dbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(CceTableName),
		Item: map[string]types.AttributeValue{
			"ID":   &types.AttributeValueMemberS{Value: cce.ID},
			"Name": &types.AttributeValueMemberS{Value: cce.Name},
		},
	})
	if err != nil {
		http.Error(w, "Failed to add CCE", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("CCE added successfully"))
}

// UpdateCCE - Update CCE by ID
func (h *CCEHandler) UpdateCCE(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cceID := vars["id"]

	var newCCE models.CCE
	err := json.NewDecoder(r.Body).Decode(&newCCE)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := dbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(CceTableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: cceID},
		},
	})
	if err != nil || result.Item == nil {
		http.Error(w, "CCE not found", http.StatusNotFound)
		return
	}

	var existingCCE models.CCE
	err = attributevalue.UnmarshalMap(result.Item, &existingCCE)
	if err != nil {
		http.Error(w, "Failed to unmarshal CCE data", http.StatusInternalServerError)
		return
	}

	// Update fields
	if newCCE.Name != "" {
		existingCCE.Name = newCCE.Name
	}

	_, err = dbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(CceTableName),
		Item: map[string]types.AttributeValue{
			"ID":   &types.AttributeValueMemberS{Value: existingCCE.ID},
			"Name": &types.AttributeValueMemberS{Value: existingCCE.Name},
		},
	})
	if err != nil {
		http.Error(w, "Failed to update CCE", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("CCE updated successfully"))
}

// DeleteCCE - Delete CCE by ID
func (h *CCEHandler) DeleteCCE(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cceID := vars["id"]

	_, err := dbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(CceTableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: cceID},
		},
	})
	if err != nil {
		http.Error(w, "Failed to delete CCE", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("CCE deleted successfully"))
}
