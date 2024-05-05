package notification

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"github.com/primexz/KrakenDCA/logger"
)

var log *logger.Logger

func init() {
	log = logger.NewLogger("gotify")
}

type GotifyMessage struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
}

func SendPushNotification(title string, messageText string) {
	appToken := os.Getenv("GOTIFY_APP_TOKEN")
	rawUrl := os.Getenv("GOTIFY_URL")

	if appToken == "" || rawUrl == "" {
		log.Debug("Skip sending push notification, not configured")
		return
	}

	reqUrl := url.URL{
		Host:   rawUrl,
		Path:   "/message",
		Scheme: "https",
	}

	q := reqUrl.Query()
	q.Set("token", appToken)
	reqUrl.RawQuery = q.Encode()

	message := &GotifyMessage{
		Title:    title,
		Message:  messageText,
		Priority: 5,
	}

	bodyBytes, err := json.Marshal(message)
	if err != nil {
		log.Error(err)
	}

	reader := bytes.NewReader(bodyBytes)
	rsp, err := http.Post(reqUrl.String(), "application/json", reader)

	if err != nil {
		log.Error("failed to send push notification", err)
	}

	if rsp.StatusCode != 200 {
		log.Error("failed to send push notification, status code", rsp.StatusCode)
	}
}
