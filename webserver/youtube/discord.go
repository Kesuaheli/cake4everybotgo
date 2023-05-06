// Copyright 2023 Kesuaheli
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package youtube

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var dcSession *discordgo.Session
var dcHandler func(*discordgo.Session, *Video)
var subscribtions = make(map[string]bool)

// SetDiscordSession sets the discord.Sesstion to use for calling
// event handlers.
func SetDiscordSession(s *discordgo.Session) {
	dcSession = s
}

// SetDiscordHandler sets the function to use when calling event
// handlers.
func SetDiscordHandler(f func(*discordgo.Session, *Video)) {
	dcHandler = f
}

// SubscribeChannel subscribe to the event listener for new videos of
// the given channel id.
func SubscribeChannel(channelID string) {
	if !subscribtions[channelID] {
		subscribtions[channelID] = true
		log.Printf("YouTube: subscribed '%s' for announcements", channelID)
	}
}

// UnsubscribeChannel removes the given channel id from the
// subscription list and no longer sends events.
func UnsubscribeChannel(channelID string) {
	if subscribtions[channelID] {
		delete(subscribtions, channelID)
		log.Printf("YouTube: unsubscribed '%s' from announcements", channelID)
	}
}
