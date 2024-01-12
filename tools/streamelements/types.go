package streamelements

import (
	"net/http"
	"time"
)

// Streamelements is the base type for communication with the streamelements API.
type Streamelements struct {
	c     *http.Client
	token string
}

// CurrentUserChannel represents the return type for the '/channels/me' endpoint.
type CurrentUserChannel struct {
	ID             string      `json:"_id"`
	Username       string      `json:"username"`
	AvatarURL      string      `json:"avatar"`
	Channels       []*Channel1 `json:"channels"`
	PrimaryChannel string      `json:"primaryChannel"`
	Teams          []string    `json:"teams"`
	LastLogin      time.Time   `json:"lastLogin"`
	CreatedAt      time.Time   `json:"createdAt"`
	UpdatedAt      time.Time   `json:"updatedAt"`
	Suspended      bool        `json:"suspended"`
}

// Channel1 represents the return type for the '/users/channels' endpoint.
type Channel1 struct {
	SimpleChannelDetails
	EmailAddress string       `json:"email"`
	Type         string       `json:"type"`
	Role         string       `json:"role"`
	Country      string       `json:"country"`
	Moderators   []*Moderator `json:"moderators"`
	LastLogin    time.Time    `json:"lastLogin"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
	Suspended    bool         `json:"suspended"`
}

// SimpleChannelDetails represents the return type for the '/channels/{channel}' endpoint.
type SimpleChannelDetails struct {
	ID              string  `json:"_id"`
	Profile         Profile `json:"profile"`
	Provider        string  `json:"provider"`
	ProviderID      string  `json:"providerId"`
	Username        string  `json:"username"`
	Alias           string  `json:"alias"`
	DisplayName     string  `json:"displayName"`
	AvatarURL       string  `json:"avatar"`
	BroadcasterType string  `json:"broadcasterType"`
	Inactive        bool    `json:"inactive"`
	IsPartner       bool    `json:"isPartner"`
}

// ChannelDetails represents the return type for the '/channels/{channel}/details' endpoint.
type ChannelDetails struct {
	SimpleChannelDetails
	Email          string    `json:"email"`
	ProviderEmails []string  `json:"providerEmails"`
	Users          []User    `json:"users"`
	Country        string    `json:"country"`
	LastJWTToken   string    `json:"lastJWTToken"`
	AccessToken    string    `json:"accessToken,omitempty"`
	APIToken       string    `json:"apiToken"`
	AB             []string  `json:"ab"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	LastLogin      time.Time `json:"lastLogin"`
	Verified       bool      `json:"verified"`
	Suspended      bool      `json:"suspended"`
	NullChannel    bool      `json:"nullChannel"`
}

// Profile represents a streamelements user profile.
type Profile struct {
	Title          string `json:"title"`
	HeaderImageURL string `json:"headerImage"`
}

// Moderator represents a moderator for a channel.
type Moderator struct {
	User User13 `json:"user"`
	Type string `json:"type"`
}

// User represents a streamelements user.
type User struct {
	User       string `json:"user"`
	ProviderID string `json:"providerId"`
	Role       string `json:"role"`
}

// User13 is a variation of the User type.
type User13 struct {
	ID        string `json:"_id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar"`
}

// UserPoints represents the return type of the '/points/{channel}/{user}' endpoint.
type UserPoints struct {
	ChannelID     string `json:"channel"`
	Username      string `json:"username"`
	Points        int    `json:"points"`
	PointsAlltime int    `json:"pointsAlltime"`
	Watchtime     int    `json:"watchtime"`
	Rank          int    `json:"rank"`
}
