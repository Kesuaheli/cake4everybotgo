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
	"fmt"
	logger "log"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	// Prefix for translation key, i.e.:
	//   key := tp+"base" // => adventcalendar
	tp = "discord.command.adventcalendar."
)

var log = logger.New(logger.Writer(), "[Advent] ", logger.LstdFlags|logger.Lmsgprefix)

type adventcalendarBase struct {
	util.InteractionUtil
	member *discordgo.Member
	user   *discordgo.User
}

type giveawayEntry struct {
	userID    string
	weight    int
	lastEntry time.Time
}

func (e giveawayEntry) toEmbedField(s *discordgo.Session) (f *discordgo.MessageEmbedField) {
	return &discordgo.MessageEmbedField{
		Name:   e.userID,
		Value:  fmt.Sprintf("<@%s>\n%d tickets\nlast entry: %s", e.userID, e.weight, e.lastEntry.Format(time.DateOnly)),
		Inline: true,
	}
}
