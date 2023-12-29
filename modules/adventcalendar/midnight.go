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
	"cake4everybot/database"
	"cake4everybot/util"
	"fmt"
	"slices"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Midnight is a scheduled function to run everyday at 0:00
func Midnight(s *discordgo.Session) {
	t := time.Now()
	if t.Month() != 12 || t.Day() == 1 || t.Day() > 25 {
		return
	}
	log.Printf("Summary for %s", t.Add(-1*time.Hour).Format("_2. Jan"))

	entries := database.GetGetAllGiveawayEntries("xmas")
	slices.SortFunc(entries, func(a, b database.GiveawayEntry) int {
		if a.Weight < b.Weight {
			return -1
		} else if a.Weight > b.Weight {
			return 1
		}
		if a.LastEntry.Before(b.LastEntry) {
			return -1
		} else if a.LastEntry.After(b.LastEntry) {
			return 1
		}
		return 0
	})
	slices.Reverse(entries)
	data := &discordgo.MessageSend{
		Embeds: splitEntriesToEmbeds(s, entries),
	}
	if len(data.Embeds) == 0 {
		data.Content = "Ticket Summary: *No Tickets!*"
	} else {
		data.Embeds[0].Title = "Current Tickets"
		if len(entries) > 1 {
			var totalTickets int
			for _, e := range entries {
				totalTickets += e.Weight
			}
			data.Embeds[0].Description = fmt.Sprintf("__Total: %d Tickets (%d users)__\nProbability per Ticket: %.2f%%\n%s", totalTickets, len(entries), 100.0/float64(totalTickets), data.Embeds[0].Description)
		}
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

func splitEntriesToEmbeds(s *discordgo.Session, entries []database.GiveawayEntry) []*discordgo.MessageEmbed {
	var totalTickets int
	for _, e := range entries {
		totalTickets += e.Weight
	}
	numEmbeds := len(entries)/25 + 1
	embeds := make([]*discordgo.MessageEmbed, 0, numEmbeds)
	for i, e := range entries {
		if i%25 == 0 {
			new := &discordgo.MessageEmbed{}
			if numEmbeds > 1 {
				new.Description = fmt.Sprintf("Page %d/%d", i/25+1, numEmbeds)
			}
			util.SetEmbedFooter(s, "module.adventcalendar.embed_footer", new)
			embeds = append(embeds, new)
		}

		embeds[len(embeds)-1].Fields = append(embeds[len(embeds)-1].Fields, e.ToEmbedField(s, totalTickets))
	}

	return embeds
}
