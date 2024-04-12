package notification

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
)

type GotifyMessage struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
}

func SendPushNotification(title string, messageText string) {
	appToken := os.Getenv("GOTIFY_APP_TOKEN")
	rawUrl := os.Getenv("GOTIFY_URL")

	if appToken == "" || rawUrl == "" {
		log.Println("Skip sending push notification, not configured")
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
		log.Println(err)
	}

	reader := bytes.NewReader(bodyBytes)
	rsp, err := http.Post(reqUrl.String(), "application/json", reader)

	if err != nil {
		log.Println("failed to send push notification", err)
	}

	if rsp.StatusCode != 200 {
		log.Println("failed to send push notification, status code", rsp.StatusCode)
	}
}
