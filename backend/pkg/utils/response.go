package utils

import (
	"encoding/json"
	"net/http"

	"backend/pkg/errors"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func RespondWithError(w http.ResponseWriter, err error) {
	switch err {
	case errors.ErrNotFound:
		RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	case errors.ErrInvalidInput:
		RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.ErrUnauthorized:
		RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
	default:
		RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": errors.ErrInternal.Error()})
	}
}
