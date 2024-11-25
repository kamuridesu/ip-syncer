package client

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/kamuridesu/ip-syncer/internal/shared"
)

func SendRequest() error {
	nameReader := strings.NewReader(shared.Info.Name)
	req, err := http.NewRequest("POST", shared.Info.ServerEndpoint, nameReader)
	if err != nil {
		slog.Error("Failed to create request object")
	}
	req.Header.Set("Authorization", shared.Info.AuthKey)

	client := &http.Client{}
	slog.Info("Sending request")
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Failed to send request")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Failed to send request")
		return err
	}

	slog.Info("Request sent successfully")

	return nil
}

func Start() {
	for {
		go func() {
			err := SendRequest()
			if err != nil {
				slog.Error("Failed to send request")
			}
		}()
		time.Sleep(5 * time.Second)
	}
}
