package api

import (
	"encoding/json"
	"net/http"
)

type BaseHandler struct {
}

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{}
}


func (h *BaseHandler) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   nil,
		"data":    nil,
		"success": true,
	})
}

func (h *BaseHandler) Success(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   nil,
		"data":    data,
		"success": true,
	})
}

func (h *BaseHandler) BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	h.errorHandler(w, r, http.StatusBadRequest, err)
}

func (h *BaseHandler) Internal(w http.ResponseWriter, r *http.Request, err error) {
	h.errorHandler(w, r, http.StatusInternalServerError, err)
}

func (h *BaseHandler) errorHandler(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   err,
		"data":    nil,
		"success": false,
	})
}
