package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/kamuridesu/gomechan/core/response"
	db "github.com/kamuridesu/ip-syncer/internal/database"
	"github.com/kamuridesu/ip-syncer/internal/shared"
)

type IHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	Database db.Database
}

func NewHandler(dbType string, info string) (*Handler, error) {
	database, err := db.New(dbType, info)
	if err != nil {
		return nil, err
	}
	return &Handler{Database: database}, nil
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	responseW := response.New(&w, r)

	requestData, err := io.ReadAll(r.Body)
	if err != nil {
		responseW.AsJson(http.StatusBadRequest, map[string]any{"error": "Failed to read request body"})
		return
	}

	ipInfo := &shared.IPInfo{}
	err = json.Unmarshal(requestData, &ipInfo)
	if err != nil {
		responseW.AsJson(http.StatusBadRequest, map[string]any{"error": "Failed to unmarshal request body"})
		return
	}

}
