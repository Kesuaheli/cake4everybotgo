package twitch

import (
	"time"
)

// Subscription represents a single subscription to an event.
type Subscription struct {
	// ID is the unique identifier for this subscription.
	ID string `json:"id"`

	// Status is the status of this subscription e.g. it is set to "enabled" when successfully
	// verified and active.
	Status string `json:"status"`

	// Type is the acual event that triggers.
	Type string `json:"type"`

	// Version is the version number of the event defined in the Type field.
	Version string `json:"version"`

	// Condition contains a list of key-value-pairs of conditions. The 'key' gives a variable to check
	// and 'value' the value to match. For example "broadcaster_user_id":"12345" requires the event
	// to be triggered at the channel of user "12345".
	Condition map[string]string `json:"condition"`

	// Transport gives information about how this subscription is (or will be) delivered.
	Transport SubscriptionTransport `json:"transport"`

	// CreatedAt is the timestamp of creation of this subscription.
	CreatedAt time.Time `json:"created_at"`

	// The amount points this subscription costs. The cost is added to a global count. Each
	// application has a fixed amount of available points to use.
	Cost int `json:"cost"`
}

// SubscriptionTransport gives information about how a subscription is (or will be) delivered.
type SubscriptionTransport struct {
	// Method is either set to "webhook" or "websocket".
	Method string `json:"method"`

	// WebhookCallbackURI gives the complete URI of the webhook.
	//
	// Only when Method == "webhook"
	WebhookCallbackURI string `json:"callback,omitempty"`

	// WebhookSecret is the secret given with the creation of the subscription to veryfiy its
	// correctness.
	//
	// Only when Method == "webhook"
	WebhookSecret string `json:"secret,omitempty"`

	// WebSocketSessionID is the ID the welcome message returns, when connecting to the twitch
	// websocket. More information needed.
	//
	// Only when Method == "websocket"
	WebSocketSessionID string `json:"session_id,omitempty"`
}
