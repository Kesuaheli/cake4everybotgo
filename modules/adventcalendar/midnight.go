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

package adventcalendar

import (
	"cake4everybot/util"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Midnight is a scheduled function to run everyday at 0:00
func Midnight(s *discordgo.Session) {
	t := time.Now()
	if t.Month() != 12 || t.Day() > 24 {
		return
	}
	log.Printf("New Post for %s", t.Format("_2. Jan"))

	var fields []*discordgo.MessageEmbedField
	for _, e := range getGetAllEntries() {
		fields = append(fields, e.toEmbedField(s))
	}

	data := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{{
			Title:  "Current Tickets",
			Fields: fields,
		}},
	}

	channels, err := util.GetChannelsFromDatabase(s, "log_channel")
	if err != nil {
		log.Printf("ERROR: Could not get advent calendar channel: %+v", err)
		return
	}

	for _, channelID := range channels {
		_, err = s.ChannelMessageSendComplex(channelID, data)
		if err != nil {
			log.Printf("ERROR: could not send log message to channel '%s': %+v", channelID, err)
			continue
		}
	}
}
