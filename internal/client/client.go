package client

import (
	"fmt"
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
		erro := fmt.Sprintf("Failed to send request: %s\n", err)
		slog.Error(erro)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error(fmt.Sprintf("Failed to send request, status is: %d\n", resp.StatusCode))
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
				erro := fmt.Sprintf("Failed to send request: %s\n", err)
				slog.Error(erro)
			}
		}()
		time.Sleep(30 * time.Second)
	}
}
