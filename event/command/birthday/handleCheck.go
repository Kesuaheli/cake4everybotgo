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
	"time"

	"cake4everybot/data/lang"
	"cake4everybot/database"
	"cake4everybot/event/command/util"

	"github.com/bwmarrin/discordgo"
)

// Check checks if there are any birthdays on the current date
// (time.Now()), if so announce them in the desired channel.
func Check(s *discordgo.Session) {
	var guildID, channelID uint64
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
	e, n := birthdayAnnounceEmbed(s, birthdays)
	if n <= 0 {
		return
	}

	for rows.Next() {
		err = rows.Scan(&guildID, &channelID)
		if err != nil {
			log.Printf("Error on scanning birthday channel ID from database %v\n", err)
			continue
		}

		channel, err := s.Channel(fmt.Sprint(channelID))
		if err != nil {
			log.Printf("Error on getting birthday channel for id: %v\n", err)
			return
		}
		if channel.GuildID != fmt.Sprint(guildID) {
			log.Printf("Warning: tried to announce birthdays in channel/%d/%d, but this channel is from guild: '%s'\n", guildID, channelID, channel.GuildID)
			return
		}

		// announce
		_, err = s.ChannelMessageSendEmbed(channel.ID, e)
		if err != nil {
			log.Printf("Error on sending todays birthday announcement: %s\n", err)
		}
	}
}

// birthdayAnnounceEmbed returns the embed, that contains all
// birthdays and 'n' as the number of birthdays, which is always
// len(b)
func birthdayAnnounceEmbed(s *discordgo.Session, b []birthdayEntry) (e *discordgo.MessageEmbed, n int) {
	var title, fValue string

	switch len(b) {
	case 0:
		title = lang.Get(tp+"msg.announce.0", lang.FallbackLang())
	case 1:
		title = lang.Get(tp+"msg.announce.1", lang.FallbackLang())
	default:
		format := lang.Get(tp+"msg.announce", lang.FallbackLang())
		title = fmt.Sprintf(format, fmt.Sprint(len(b)))
	}

	for _, b := range b {
		mention := fmt.Sprintf("<@%d>", b.ID)

		if b.Year == 0 {
			fValue += fmt.Sprintf("%s\n", mention)
		} else {
			format := lang.Get(tp+"msg.announce.with_age", lang.FallbackLang())
			format += "\n"
			fValue += fmt.Sprintf(format, mention, fmt.Sprint(b.Age()+1))
		}
	}

	e = &discordgo.MessageEmbed{
		Title: title,
		Color: 0xFFD700,
	}

	if len(b) == 0 {
		e.Color = 0xFF0000
		e.Description = lang.Get(tp+"msg.announce.0.description", lang.FallbackLang())
	} else {
		e.Color = 0xFFD700
		e.Fields = []*discordgo.MessageEmbedField{{
			Name:  lang.Get(tp+"msg.announce.congratulate", lang.FallbackLang()),
			Value: fValue,
		}}
	}

	util.SetEmbedFooter(s, tp+"display", e)

	return e, len(b)
}
