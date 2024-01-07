package twitch

import (
	"time"
)

type Subscription struct {
	ID        string                `json:"id"`
	Status    string                `json:"status"`
	Type      string                `json:"type"`
	Version   string                `json:"version"`
	Condition map[string]string     `json:"condition"`
	Transport SubscriptionTransport `json:"transport"`
	CreatedAt time.Time             `json:"created_at"`
	Cost      int                   `json:"cost"`
}
type SubscriptionTransport struct {
	Method             string `json:"method"`
	WebhookCallbackURI string `json:"callback,omitempty"`
	WebhookSecret      string `json:"secret,omitempty"`
	WebSocketSessionID string `json:"session_id,omitempty"`
}
