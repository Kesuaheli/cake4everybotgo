package twitch

import "github.com/bwmarrin/discordgo"

var dcSession *discordgo.Session
var dcHandler func(*discordgo.Session, *RawEvent)
var subscribtions = make(map[string]bool)

// SetDiscordSession sets the discord.Sesstion to use for calling
// event handlers.
func SetDiscordSession(s *discordgo.Session) {
	dcSession = s
}

// SetDiscordHandler sets the function to use when calling event
// handlers.
func SetDiscordHandler(f func(*discordgo.Session, *RawEvent)) {
	dcHandler = f
}

// SubscribeChannel subscribe to the event listener for new videos of
// the given channel id.
func SubscribeChannel(channelID string) {
	if !subscribtions[channelID] {
		subscribtions[channelID] = true
		log.Printf("subscribed '%s' for announcements", channelID)
	}
}

// UnsubscribeChannel removes the given channel id from the
// subscription list and no longer sends events.
func UnsubscribeChannel(channelID string) {
	if subscribtions[channelID] {
		delete(subscribtions, channelID)
		log.Printf("unsubscribed '%s' from announcements", channelID)
	}
}
