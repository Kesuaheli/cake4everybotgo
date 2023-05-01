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
	"github.com/bwmarrin/discordgo"
)

var dcSession *discordgo.Session
var subscribtionMap map[string]func(*discordgo.Session, Feed) = make(map[string]func(*discordgo.Session, Feed))

// SetDiscordSession sets the discord.Sesstion to use for calling
// event handlers.
func SetDiscordSession(s *discordgo.Session) {
	dcSession = s
}

// SubscribeChannel subscribe to the event listener for new videos of
// the given channel id.
func SubscribeChannel(channelID string, f func(s *discordgo.Session, e Feed)) {
	subscribtionMap[channelID] = f
}
