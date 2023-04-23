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

package birthday

import (
	"fmt"
	"log"
	"strings"
	"time"

	"cake4everybot/database"

	"github.com/bwmarrin/discordgo"
)

func Check(s *discordgo.Session) {
	var guild_id, channel_id uint64
	rows, err := database.Query("SELECT id,birthday_id FROM guilds")
	if err != nil {
		log.Printf("Error on getting birthday channel IDs from database: %v\n", err)
	}
	defer rows.Close()

	now := time.Now()
	birthdays, err := getBirthdaysDate(now.Day(), int(now.Month()))
	if err != nil {
		log.Printf("Error on getting todays birthdays from database: %v\n", err)
	}
	if len(birthdays) == 0 {
		return
	}

	for rows.Next() {
		err = rows.Scan(&guild_id, &channel_id)
		if err != nil {
			log.Printf("Error on scanning birthday channel ID from database %v\n", err)
			continue
		}

		channel, err := s.Channel(fmt.Sprint(channel_id))
		if err != nil {
			log.Printf("Error on getting birthday channel for id: %v\n", err)
			return
		}
		if channel.GuildID != fmt.Sprint(guild_id) {
			log.Printf("Warning: tried to announce birthdays in channel/%d/%d, but this channel is from guild: '%s'\n", guild_id, channel_id, channel.GuildID)
			return
		}

		announceBirthdays(s, channel, birthdays)
	}
}

func announceBirthdays(s *discordgo.Session, channel *discordgo.Channel, birthdays []birthdayEntry) {
	var (
		n   int
		msg string
	)

	for _, b := range birthdays {
		m, err := s.GuildMember(channel.GuildID, fmt.Sprint(b.ID))
		if err != nil {
			if !strings.HasPrefix(err.Error(), "HTTP 404 Not Found") {
				log.Printf("Error on get guild member '%d' in guild '%s': %v\n", b.ID, channel.GuildID, err)
			}
			continue
		}

		n = n + 1
		var age string
		if b.Year > 0 {
			age = fmt.Sprintf(" turns %d", time.Now().Year()-b.Year)
		}
		msg = fmt.Sprintf("%s\n%s%s!", msg, m.Mention(), age)
	}

	if n == 0 {
		return
	} else if n == 1 {
		msg = fmt.Sprintf("Theres a birthday today!%s", msg)
	} else {
		msg = fmt.Sprintf("There are %d birthdays today!%s", n, msg)
	}

	_, err := s.ChannelMessageSend(channel.ID, msg)
	if err != nil {
		log.Printf("Error on sending todays birthday message: %s\n", err)
	}

}
