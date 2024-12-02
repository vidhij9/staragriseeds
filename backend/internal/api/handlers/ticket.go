package handlers

import (
	"backend/internal/models"
	"backend/internal/service"
	"backend/pkg/errors"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gorilla/mux"
)

const TicketTableName = "Tickets"

type TicketHandler struct {
	ticketService *service.TicketService
	farmerService *service.FarmerService
}

func NewTicketHandler(ticketService *service.TicketService, farmerService *service.FarmerService) *TicketHandler {
	return &TicketHandler{
		ticketService: ticketService,
		farmerService: farmerService,
	}
}

// GetTicket - Retrieve ticket by ID
func (h *TicketHandler) GetTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ticketID := vars["id"]

	result, err := dbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(TicketTableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: ticketID},
		},
	})
	if err != nil || result.Item == nil {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	var ticket models.Ticket
	err = attributevalue.UnmarshalMap(result.Item, &ticket)
	if err != nil {
		http.Error(w, "Failed to unmarshal ticket data", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ticket)
}

// GetTickets - Retrieve all tickets
func (h *TicketHandler) GetTickets(w http.ResponseWriter, r *http.Request) {
	result, err := dbClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(TicketTableName),
	})
	if err != nil {
		http.Error(w, "Failed to fetch tickets", http.StatusInternalServerError)
		return
	}

	var tickets []models.Ticket
	err = attributevalue.UnmarshalListOfMaps(result.Items, &tickets)
	if err != nil {
		http.Error(w, "Failed to unmarshal tickets data", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tickets)
}

// GetTicketsByFarmerContact - Retrieve all tickets by a farmer's contact
func (h *TicketHandler) GetTicketsByFarmer(w http.ResponseWriter, r *http.Request) {
	contact := r.URL.Query().Get("contact")
	if contact == "" {
		errors.WriteJSONError(w, http.StatusBadRequest, "Contact is required")
		return
	}

	farmer, err := h.farmerService.GetFarmerByContact(r.Context(), contact)
	if err != nil {
		errors.WriteJSONError(w, http.StatusInternalServerError, "Failed to get farmer")
		return
	}

	tickets, err := h.ticketService.GetTicketsByFarmerContact(r.Context(), farmer)
	if err != nil {
		errors.WriteJSONError(w, http.StatusInternalServerError, "Failed to get tickets")
		return
	}

	json.NewEncoder(w).Encode(tickets)
}

// GetTicketsByCCE - Retrieve all tickets by a CCE's contact
func (h *TicketHandler) GetTicketsByCCE(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cceID := vars["cceId"]
	if cceID == "" {
		errors.WriteJSONError(w, http.StatusBadRequest, "CCE ID is required")
		return
	}

	tickets, err := h.ticketService.GetTicketsByCCE(r.Context(), cceID)
	if err != nil {
		errors.WriteJSONError(w, http.StatusInternalServerError, "Failed to get tickets")
		return
	}

	json.NewEncoder(w).Encode(tickets)
}

func (h *TicketHandler) GetTicketsByCCEWithDateFilter(w http.ResponseWriter, r *http.Request) {
	cceID := r.URL.Query().Get("cceId")
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	if cceID == "" || startDate == "" || endDate == "" {
		errors.WriteJSONError(w, http.StatusBadRequest, "CCE ID, start date, and end date are required")
		return
	}

	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		errors.WriteJSONError(w, http.StatusBadRequest, "Invalid start date format")
		return
	}

	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		errors.WriteJSONError(w, http.StatusBadRequest, "Invalid end date format")
		return
	}

	tickets, err := h.ticketService.GetTicketsByCCEWithDateFilter(r.Context(), cceID, start, end)
	if err != nil {
		errors.WriteJSONError(w, http.StatusInternalServerError, "Failed to get tickets")
		return
	}

	json.NewEncoder(w).Encode(tickets)
}

func (h *TicketHandler) GetTicketsByCCEAndStatus(w http.ResponseWriter, r *http.Request) {
	cceID := r.URL.Query().Get("cceId")
	status := r.URL.Query().Get("status")

	if cceID == "" || status == "" {
		errors.WriteJSONError(w, http.StatusBadRequest, "CCE ID and status are required")
		return
	}

	tickets, err := h.ticketService.GetTicketsByCCEAndStatus(r.Context(), cceID, status)
	if err != nil {
		errors.WriteJSONError(w, http.StatusInternalServerError, "Failed to get tickets")
		return
	}

	json.NewEncoder(w).Encode(tickets)
}

func (h *TicketHandler) GetTicketsWithStatusAndSort(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	sortBy := r.URL.Query().Get("sortBy")
	sortOrder := r.URL.Query().Get("sortOrder")

	if status == "" {
		errors.WriteJSONError(w, http.StatusBadRequest, "Status is required")
		return
	}

	tickets, err := h.ticketService.GetTicketsWithStatusAndSort(r.Context(), status, sortBy, sortOrder)
	if err != nil {
		errors.WriteJSONError(w, http.StatusInternalServerError, "Failed to get tickets")
		return
	}

	json.NewEncoder(w).Encode(tickets)
}

// CreateTicket - Add new ticket
func (h *TicketHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	var ticket models.Ticket
	err := json.NewDecoder(r.Body).Decode(&ticket)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err = dbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(TicketTableName),
		Item: map[string]types.AttributeValue{
			"ID":       &types.AttributeValueMemberS{Value: ticket.ID},
			"FarmerID": &types.AttributeValueMemberS{Value: ticket.FarmerID},
			"CCEID":    &types.AttributeValueMemberS{Value: ticket.CCEID},
			"Status":   &types.AttributeValueMemberS{Value: ticket.Status},
		},
	})
	if err != nil {
		http.Error(w, "Failed to add ticket", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Ticket added successfully"))
}

// UpdateTicket - Update ticket by ID
func (h *TicketHandler) UpdateTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ticketID := vars["id"]

	var newTicket models.Ticket
	err := json.NewDecoder(r.Body).Decode(&newTicket)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := dbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(TicketTableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: ticketID},
		},
	})
	if err != nil || result.Item == nil {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	var existingTicket models.Ticket
	err = attributevalue.UnmarshalMap(result.Item, &existingTicket)
	if err != nil {
		http.Error(w, "Failed to unmarshal ticket data", http.StatusInternalServerError)
		return
	}

	// Update fields
	if newTicket.FarmerID != "" {
		existingTicket.FarmerID = newTicket.FarmerID
	}
	if newTicket.CCEID != "" {
		existingTicket.CCEID = newTicket.CCEID
	}
	if newTicket.Status != "" {
		existingTicket.Status = newTicket.Status
	}

	_, err = dbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(TicketTableName),
		Item: map[string]types.AttributeValue{
			"ID":       &types.AttributeValueMemberS{Value: existingTicket.ID},
			"FarmerID": &types.AttributeValueMemberS{Value: existingTicket.FarmerID},
			"CceID":    &types.AttributeValueMemberS{Value: existingTicket.CCEID},
			"Status":   &types.AttributeValueMemberS{Value: existingTicket.Status},
		},
	})
	if err != nil {
		http.Error(w, "Failed to update ticket", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Ticket updated successfully"))
}

// DeleteTicket - Delete ticket by ID
func (h *TicketHandler) DeleteTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ticketID := vars["id"]

	_, err := dbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(TicketTableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: ticketID},
		},
	})
	if err != nil {
		http.Error(w, "Failed to delete ticket", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Ticket deleted successfully"))
}
