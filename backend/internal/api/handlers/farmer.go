package handlers

import (
	"backend/internal/models"
	"backend/internal/service"
	"backend/pkg/errors"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gorilla/mux"
)

const FarmerTableName = "Farmers"

type FarmerHandler struct {
	farmerService *service.FarmerService
}

func NewFarmerHandler(farmerService *service.FarmerService) *FarmerHandler {
	return &FarmerHandler{farmerService: farmerService}
}

// CreateFarmer - Add a new farmer
func (h *FarmerHandler) CreateFarmer(w http.ResponseWriter, r *http.Request) {
	var farmer models.Farmer
	err := json.NewDecoder(r.Body).Decode(&farmer)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err = dbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(FarmerTableName),
		Item: map[string]types.AttributeValue{
			"ID":       &types.AttributeValueMemberS{Value: farmer.ID},
			"Name":     &types.AttributeValueMemberS{Value: farmer.Name},
			"Contact":  &types.AttributeValueMemberS{Value: farmer.Contact},
			"State":    &types.AttributeValueMemberS{Value: farmer.State},
			"District": &types.AttributeValueMemberS{Value: farmer.District},
			"Tehsil":   &types.AttributeValueMemberS{Value: farmer.Tehsil},
			"Village":  &types.AttributeValueMemberS{Value: farmer.Village},
			"Pincode":  &types.AttributeValueMemberS{Value: farmer.Pincode},
			"Address":  &types.AttributeValueMemberS{Value: farmer.Address},
			"Tag":      &types.AttributeValueMemberS{Value: farmer.Tag},
			"Crop":     &types.AttributeValueMemberSS{Value: farmer.Crop}, // Use string set for crops
		},
	})

	if err != nil {
		http.Error(w, "Failed to add farmer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Farmer added successfully")
}

// UpdateFarmer - Update farmer by ID
func (h *FarmerHandler) UpdateFarmer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	farmerID := vars["id"]

	var newFarmer models.Farmer
	err := json.NewDecoder(r.Body).Decode(&newFarmer)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Retrieve the existing farmer
	result, err := dbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(FarmerTableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: farmerID},
		},
	})
	if err != nil || result.Item == nil {
		http.Error(w, "Farmer not found", http.StatusNotFound)
		return
	}

	var existingFarmer models.Farmer
	err = attributevalue.UnmarshalMap(result.Item, &existingFarmer)
	if err != nil {
		http.Error(w, "Failed to unmarshal farmer data", http.StatusInternalServerError)
		return
	}

	// Only update fields that are provided in the request body
	if newFarmer.Name != "" {
		existingFarmer.Name = newFarmer.Name
	}
	if newFarmer.Contact != "" {
		existingFarmer.Contact = newFarmer.Contact
	}
	if newFarmer.State != "" {
		existingFarmer.State = newFarmer.State
	}
	if newFarmer.District != "" {
		existingFarmer.District = newFarmer.District
	}
	if newFarmer.Tehsil != "" {
		existingFarmer.Tehsil = newFarmer.Tehsil
	}
	if newFarmer.Village != "" {
		existingFarmer.Village = newFarmer.Village
	}
	if newFarmer.Pincode != "" {
		existingFarmer.Pincode = newFarmer.Pincode
	}
	if newFarmer.Address != "" {
		existingFarmer.Address = newFarmer.Address
	}
	if newFarmer.Tag != "" {
		existingFarmer.Tag = newFarmer.Tag
	}
	if len(newFarmer.Crop) > 0 {
		existingFarmer.Crop = newFarmer.Crop
	}

	// Update the farmer in DynamoDB
	_, err = dbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(FarmerTableName),
		Item: map[string]types.AttributeValue{
			"ID":       &types.AttributeValueMemberS{Value: existingFarmer.ID},
			"Name":     &types.AttributeValueMemberS{Value: existingFarmer.Name},
			"Contact":  &types.AttributeValueMemberS{Value: existingFarmer.Contact},
			"State":    &types.AttributeValueMemberS{Value: existingFarmer.State},
			"District": &types.AttributeValueMemberS{Value: existingFarmer.District},
			"Tehsil":   &types.AttributeValueMemberS{Value: existingFarmer.Tehsil},
			"Village":  &types.AttributeValueMemberS{Value: existingFarmer.Village},
			"Pincode":  &types.AttributeValueMemberS{Value: existingFarmer.Pincode},
			"Address":  &types.AttributeValueMemberS{Value: existingFarmer.Address},
			"Tag":      &types.AttributeValueMemberS{Value: existingFarmer.Tag},
			"Crop":     &types.AttributeValueMemberSS{Value: existingFarmer.Crop},
		},
	})
	if err != nil {
		http.Error(w, "Failed to update farmer", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Farmer updated successfully")
}

// GetFarmer - Retrieve farmer by ID
func (h *FarmerHandler) GetFarmer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	farmerID := vars["id"]

	result, err := dbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(FarmerTableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: farmerID},
		},
	})

	if err != nil || result.Item == nil {
		http.Error(w, "Farmer not found", http.StatusNotFound)
		return
	}

	var farmer models.Farmer
	err = attributevalue.UnmarshalMap(result.Item, &farmer)
	if err != nil {
		http.Error(w, "Failed to unmarshal farmer data", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(farmer)
}

// GetFarmerByContact - Retrieve farmer by contact number
func (h *FarmerHandler) GetFarmerByContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contact := vars["contact"]

	if contact == "" {
		errors.WriteJSONError(w, http.StatusBadRequest, "Contact is required")
		return
	}

	farmer, err := h.farmerService.GetFarmerByContact(r.Context(), contact)
	if err != nil {
		errors.WriteJSONError(w, http.StatusInternalServerError, "Failed to get farmer")
		return
	}

	json.NewEncoder(w).Encode(farmer)
}

// GetFarmers - Retieve farmers by filters{crop, village, etc}
func (h *FarmerHandler) GetFarmers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	queryParams := r.URL.Query()
	filters := map[string]string{
		"crop":     queryParams.Get("crop"),
		"district": queryParams.Get("district"),
		"village":  queryParams.Get("village"),
		"pincode":  queryParams.Get("pincode"),
	}

	farmers, err := h.farmerService.ListFarmersWithFilters(r.Context(), filters)
	if err != nil {
		errors.WriteJSONError(w, http.StatusInternalServerError, "Failed to list farmers:"+err.Error())
		return
	}

	json.NewEncoder(w).Encode(farmers)
}

// DeleteFarmer - Delete a farmer by ID
func (h *FarmerHandler) DeleteFarmer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	farmerID := vars["id"]

	_, err := dbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(FarmerTableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: farmerID},
		},
	})
	if err != nil {
		http.Error(w, "Failed to delete farmer", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Farmer deleted successfully")
}
