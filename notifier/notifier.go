package notifier

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/beck-8/filecoin-checker/config"

	"github.com/rs/zerolog/log"
)

// NotifyRequest defines the structure for sending notification requests
// apprise format
type NotifyRequest struct {
	URLs  string `json:"urls"`  // Notification target URLs (like mailto://, discord://)
	Body  string `json:"body"`  // Notification content
	Title string `json:"title"` // Notification title
}

func SendNotify(miner, body, title string, recipientURLs []string, serverURL string) error {
	if serverURL == "" {
		if config.Global.Global.AppriseAPIServer != "" {
			serverURL = config.Global.Global.AppriseAPIServer
		} else {
			log.Warn().Str("miner", miner).Msg("Notification server address not configured")
			return nil
		}
	}

	if len(recipientURLs) == 0 {
		if len(config.Global.Global.RecipientURLs) != 0 {
			recipientURLs = config.Global.Global.RecipientURLs
		} else {
			log.Warn().Str("miner", miner).Msg("Notification targets not configured")
			return nil
		}
	}

	request := NotifyRequest{
		URLs:  strings.Join(recipientURLs, ","),
		Body:  body,
		Title: title,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Error().Str("miner", miner).Err(err).Msg("Failed to build request body")
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Post(serverURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error().Str("miner", miner).Err(err).Msg("Failed to send request")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Error().
			Str("miner", miner).
			Int("status_code", resp.StatusCode).
			Str("response", string(body)).
			Msg("Failed to send notification")
		return nil
	}

	log.Debug().
		Str("miner", miner).
		Str("title", title).
		Str("body", body).
		// Str("serverURL", serverURL).
		// Strs("recipientURLs", recipientURLs).
		Msg("Notification sent successfully")

	return nil
}
