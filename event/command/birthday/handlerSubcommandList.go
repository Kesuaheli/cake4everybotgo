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
	"cake4everybot/event/command/util"

	"github.com/bwmarrin/discordgo"
)

// The list subcommand. Used when executing the
// slash-command "/birthday list".
type subcommandList struct {
	Chat
	*discordgo.ApplicationCommandInteractionDataOption

	month *discordgo.ApplicationCommandInteractionDataOption // reqired
}

// Constructor for subcommandList, the struct for
// the slash-command "/birthday remove".
func (cmd Chat) subcommandList() subcommandList {
	subcommand := cmd.Interaction.ApplicationCommandData().Options[0]
	return subcommandList{
		Chat:                                    cmd,
		ApplicationCommandInteractionDataOption: subcommand,
	}
}

func (cmd subcommandList) handler() {
	for _, opt := range cmd.Options {
		switch opt.Name {
		case "month":
			cmd.month = opt
		}
	}
	month := int(cmd.month.IntValue())

	birthdays, err := cmd.getBirthdaysMonth(month)
	if err != nil {
		log.Printf("Error on get birthdays by month: %v\n", err)
		cmd.ReplyError()
		return
	}

	monthName := lang.GetSlice(tp+"month", month-1, lang.FallbackLang())
	var (
		header, key, value string
		a                  []any
	)

	switch len(birthdays) {
	case 0, 1:
		key = fmt.Sprintf("%smsg.list.total.%d", tp, len(birthdays))
		a = append(a, monthName)
	default:
		key = tp + "msg.list.total"
		a = append(a,
			fmt.Sprint(len(birthdays)),
			monthName,
		)
	}
	header = fmt.Sprintf(lang.Get(key, lang.FallbackLang()), a...)

	for _, b := range birthdays {
		var timestamp string
		if time.Until(b.Next()) <= time.Hour*24*25 {
			timestamp = fmt.Sprintf(" <t:%d:R>", b.NextUnix())
		}
		value += fmt.Sprintf("`%s` <@%d>%s\n", b.String(), b.ID, timestamp)
	}

	e := &discordgo.MessageEmbed{
		Title: fmt.Sprintf(lang.Get(tp+"msg.list", lang.FallbackLang()), monthName),
		Fields: []*discordgo.MessageEmbedField{{
			Name:   header,
			Value:  value,
			Inline: false,
		}},
		Color: 0x00FF00,
	}
	if len(birthdays) == 0 {
		e.Color = 0xFF0000
	}
	util.SetEmbedFooter(cmd.Session, tp+"display", e)

	cmd.ReplyEmbed(e)
}
