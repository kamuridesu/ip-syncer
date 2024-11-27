package server

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/kamuridesu/gomechan/core/response"
	db "github.com/kamuridesu/ip-syncer/internal/database"
	"github.com/kamuridesu/ip-syncer/internal/hosts"
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

func validateRequest(r *http.Request, responseW *response.ResponseWriter) error {
	if r.Method != http.MethodPost {
		responseW.AsJson(http.StatusMethodNotAllowed, map[string]any{"error": "Method Not Allowed"})
		return fmt.Errorf("method Not Allowed")
	}

	if r.Header.Get("Authorization") != shared.Info.AuthKey {
		responseW.AsJson(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
		return fmt.Errorf("nnauthorized")
	}
	return nil
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	responseW := response.New(&w, r)

	if err := validateRequest(r, &responseW); err != nil {
		return
	}

	hostsFile, err := hosts.ReadHostsFile()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to read hosts file: %s", err))
		responseW.AsJson(http.StatusInternalServerError, map[string]any{"error": "Internal Server Error"})
		return
	}

	ipAddr := r.Header.Get("X-Forwarded-For")
	if ipAddr == "" {
		ipAddr = strings.Split(r.RemoteAddr, ":")[0]
	}

	content, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to read request body: %s", err))
		responseW.AsJson(http.StatusInternalServerError, map[string]any{"error": "Internal Server Error"})
		return
	}

	name := string(content)
	if name == "" {
		responseW.AsJson(http.StatusBadRequest, map[string]any{"error": "Name is required"})
		return
	}
	info := shared.NewIPInfo(ipAddr, name)
	storedInfo, err := h.Database.GetByIP(ipAddr)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to get record: %s", err))
		responseW.AsJson(http.StatusInternalServerError, map[string]any{"error": "Internal Server Error"})
		return
	}

	if storedInfo != nil {
		if info.Equals(storedInfo) {
			responseW.AsJson(http.StatusOK, map[string]any{"message": "IP already exists"})
			return
		}

		hostsFile.AddOrReplaceHost(info.IP, info.Name).Save()

		err = h.Database.Update(info)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to update record: %s", err))
			responseW.AsJson(http.StatusInternalServerError, map[string]any{"error": "Internal Server Error"})
			return
		}

		responseW.AsJson(http.StatusOK, map[string]any{"message": "IP updated"})
		return
	}

	err = h.Database.Insert(info)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to insert record: %s", err))
		responseW.AsJson(http.StatusInternalServerError, map[string]any{"error": "Internal Server Error"})
		return
	}
	hostsFile.AddOrReplaceHost(info.IP, info.Name).Save()

	responseW.AsJson(http.StatusOK, map[string]any{"message": "IP added"})

}

func Start(handler *Handler) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Handle)
	slog.Info("Server started on port 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		slog.Error(fmt.Sprintf("Error starting server: %s", err))
	}
}
