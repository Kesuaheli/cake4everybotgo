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

package twitch

import (
	"log"

	twitchgo "github.com/gempir/go-twitch-irc"
	"github.com/spf13/viper"
)

// Handle adds all handler to the client
func Handle(client *twitchgo.Client) {
	client.OnConnect(func() {
		log.Printf("Twich connected as %s!\n", viper.GetString("twitch.name"))
	})

	channels := viper.GetStringSlice("twitch.channels")
	for _, channel := range channels {
		client.Join(channel)
	}
	log.Printf("Channel list set to %v\n", channels)

	client.OnNewMessage(messageHandler)

	client.OnUserJoin(func(channel, user string) {
		if user == viper.GetString("twitch.name") {
			log.Printf("Connected to %s", channel)
		} else {
			log.Printf("Twitch: %s joined %s\n", user, channel)
		}
	})
	client.OnUserPart(func(channel, user string) {
		log.Printf("Twitch: %s left %s\n", user, channel)
	})

	client.OnNewNoticeMessage(func(channel string, user twitchgo.User, message twitchgo.Message) {
		log.Printf("Twitch [Notice]: %s@%s: %s", user.DisplayName, channel, message.Raw)
	})
	client.OnNewRoomstateMessage(func(channel string, user twitchgo.User, message twitchgo.Message) {
		log.Printf("Twitch [RoomState]: %s@%s: %s", user.DisplayName, channel, message.Raw)
	})

	client.OnNewUsernoticeMessage(func(channel string, user twitchgo.User, message twitchgo.Message) {
		log.Printf("Twitch [UserNotice]: %s@%s: %s", user.DisplayName, channel, message.Raw)
	})
	client.OnNewUserstateMessage(func(channel string, user twitchgo.User, message twitchgo.Message) {
		log.Printf("Twitch [UserState]: %s@%s: %s", user.DisplayName, channel, message.Raw)
	})
}
