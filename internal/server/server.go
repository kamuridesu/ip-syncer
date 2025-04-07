package server

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/kamuridesu/gomechan/core/response"
	"github.com/kamuridesu/ip-syncer/internal/hosts"
	"github.com/kamuridesu/ip-syncer/internal/shared"
)

type Handler struct {
	hosts *hosts.Hosts
}

func NewHandler(hosts *hosts.Hosts) (*Handler, error) {

	return &Handler{hosts: hosts}, nil
}

func validateRequest(r *http.Request, responseW *response.ResponseWriter) error {
	if r.Method != http.MethodPost {
		responseW.AsJson(http.StatusMethodNotAllowed, map[string]any{"error": "Method Not Allowed"})
		return fmt.Errorf("method Not Allowed")
	}

	if r.Header.Get("Authorization") != shared.Info.AuthKey {
		responseW.AsJson(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
		return fmt.Errorf("unauthorized")
	}
	return nil
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	responseW := response.New(&w, r)

	if err := validateRequest(r, &responseW); err != nil {
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

	err = h.hosts.AddOrReplaceHost(info.IP, info.Name)
	if err != nil {
		slog.Error(fmt.Sprintf("error saving new ip, err is %s\n", err))
		responseW.AsJson(http.StatusInternalServerError, map[string]any{"error": "Internal Server Error"})
		return
	}

	responseW.AsJson(http.StatusOK, map[string]any{"message": "IP updated"})

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
