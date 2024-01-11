package twitch

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	logger "log"
	"net/http"

	"github.com/spf13/viper"
)

// RawEvent represents the http body comming with a call to the /twitch_pubsub enpoints
type RawEvent struct {
	// Challenge cointains the string to return when receiving a webhook callback verification.
	// Otherwise it is an empty string
	Challenge string `json:"challenge"`

	// Subscription contains the informations this event is about.
	Subscription Subscription `json:"subscription"`

	// Event is the actual event.
	//
	// It is not set in a webhook callback verification.
	Event interface{} `json:"event"`
}

var log = logger.New(logger.Writer(), "[WebTwitch] ", logger.LstdFlags|logger.Lmsgprefix)

// HandlePost is the HTTP/POST handler for the Twitch PubSub endpoint.
//
// It is called to handle a webhook comming from twitch. This could be a hub challenge verification
// or a event notification.
func HandlePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// before anything, check the hash
	if !verifyTwitchMessage(r.Header, body) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var rEvent RawEvent
	err = json.Unmarshal(body, &rEvent)
	if err != nil {
		log.Printf("Failed to unmarshal body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	messageType := r.Header.Get("Twitch-Eventsub-Message-Type")
	switch messageType {
	case "webhook_callback_verification":
		handleVerification(w, r, rEvent)
		return
	case "notification":
		log.Printf("Event notification: %+v", rEvent)
	default:
		log.Printf("Unknown message type '%s'", messageType)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func verifyTwitchMessage(header http.Header, body []byte) bool {
	// read more at https://dev.twitch.tv/docs/eventsub/handling-webhook-events/#verifying-the-event-message
	msgID := header.Get("Twitch-Eventsub-Message-Id")
	msgTime := header.Get("Twitch-Eventsub-Message-Timestamp")
	data := append([]byte(msgID+msgTime), body...)

	h := hmac.New(sha256.New, []byte(viper.GetString("twitch.webhookSecret")))
	h.Write(data)
	hmacData := h.Sum(nil)

	hmacHex := make([]byte, 128)
	n := hex.Encode(hmacHex, hmacData)
	hmacHex = append([]byte("sha256="), hmacHex[:n]...)

	signature := []byte(header.Get("Twitch-Eventsub-Message-Signature"))
	return hmac.Equal(hmacHex, signature)
}

func handleVerification(w http.ResponseWriter, r *http.Request, rEvent RawEvent) {
	if rEvent.Challenge == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	broadcaster := rEvent.Subscription.Condition["broadcaster_user_id"]
	if broadcaster != "404257324" {
		log.Printf("Declined verification for broadcaster '%s'!", broadcaster)
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("{\"conflict\":\"that broadcaster is not allowed\"}"))
		return
	}

	log.Printf("Accepted '%s v%s' for channel %s", rEvent.Subscription.Type, rEvent.Subscription.Version, broadcaster)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(rEvent.Challenge))
}
