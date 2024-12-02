package handlers

import (
	"backend/internal/service"
	"backend/pkg/errors"
	"encoding/json"
	"net/http"
	"time"
)

const ShootTableName = "Shoots"

type ShootHandler struct {
	shootService *service.ShootService
}

func NewShootHandler(shootService *service.ShootService) *ShootHandler {
	return &ShootHandler{shootService: shootService}
}

func (h *ShootHandler) GetAllShoots(w http.ResponseWriter, r *http.Request) {
	shootType := r.URL.Query().Get("type")
	if shootType != "" && shootType != "whatsapp" && shootType != "call" {
		errors.WriteJSONError(w, http.StatusBadRequest, "Invalid shoot type")
		return
	}

	shoots, err := h.shootService.GetAllShoots(r.Context(), shootType)
	if err != nil {
		errors.WriteJSONError(w, http.StatusInternalServerError, "Failed to get shoots")
		return
	}

	json.NewEncoder(w).Encode(shoots)
}

func (h *ShootHandler) GetShootsWithDateFilter(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	if startDate == "" || endDate == "" {
		errors.WriteJSONError(w, http.StatusBadRequest, "Start date and end date are required")
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

	shoots, err := h.shootService.GetShootsWithDateFilter(r.Context(), start, end)
	if err != nil {
		errors.WriteJSONError(w, http.StatusInternalServerError, "Failed to get shoots")
		return
	}

	json.NewEncoder(w).Encode(shoots)
}

func (h *ShootHandler) GetMissedShoots(w http.ResponseWriter, r *http.Request) {
	shoots, err := h.shootService.GetMissedShoots(r.Context())
	if err != nil {
		errors.WriteJSONError(w, http.StatusInternalServerError, "Failed to get missed shoots")
		return
	}

	json.NewEncoder(w).Encode(shoots)
}
