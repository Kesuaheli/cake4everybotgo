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
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"

	"cake4everybot/data/lang"
)

// The set subcommand. Used when executing the
// slash-command "/birthday set".
type subcommandSet struct {
	Chat
	*discordgo.ApplicationCommandInteractionDataOption

	day     *discordgo.ApplicationCommandInteractionDataOption // reqired
	month   *discordgo.ApplicationCommandInteractionDataOption // reqired
	year    *discordgo.ApplicationCommandInteractionDataOption // optional
	visible *discordgo.ApplicationCommandInteractionDataOption // optional
}

// Constructor for subcommandSet, the struct for
// the slash-command "/birthday set".
func (cmd Chat) subcommandSet() subcommandSet {
	subcommand := cmd.Interaction.ApplicationCommandData().Options[0]
	return subcommandSet{
		Chat:                                    cmd,
		ApplicationCommandInteractionDataOption: subcommand,
	}
}

func (cmd subcommandSet) handler() {
	for _, opt := range cmd.Options {
		switch opt.Name {
		case lang.GetDefault(tp + "option.set.option.day"):
			cmd.day = opt
		case lang.GetDefault(tp + "option.set.option.month"):
			cmd.month = opt
		case lang.GetDefault(tp + "option.set.option.year"):
			cmd.year = opt
		case lang.GetDefault(tp + "option.set.option.visible"):
			cmd.visible = opt
		}
	}

	switch cmd.Interaction.Type {
	case discordgo.InteractionApplicationCommandAutocomplete:
		cmd.autocompleteHandler()
	case discordgo.InteractionApplicationCommand:
		cmd.interactionHandler()
	}
}

func (cmd subcommandSet) autocompleteHandler() {
	var choices []*discordgo.ApplicationCommandOptionChoice

	if cmd.day != nil && cmd.day.Focused {
		start := cmd.day.Value.(string)

		var m int
		if cmd.month != nil {
			m = int(cmd.month.IntValue())
		}
		if m == 0 {
			m = 1
		}
		leapYear := cmd.year == nil || cmd.year.IntValue()%4 == 0

		choices = dayChoices(start, m, leapYear)
	} else if cmd.month != nil && cmd.month.Focused {
		start := cmd.month.Value.(string)

		var d int
		if cmd.day != nil {
			d = int(cmd.day.IntValue())
		}
		leapYear := cmd.year == nil || cmd.year.IntValue()%4 == 0

		locale := lang.FallbackLang()
		if user, err := cmd.Session.User(cmd.user.ID); err == nil {
			locale = user.Locale
			log.Printf("DEBUG: Users Locale '%s'\n", locale)
		}

		choices = monthChoices(start, d, leapYear)
	} else if cmd.year != nil && cmd.year.Focused {
		start := cmd.year.Value.(string)

		var d, m int
		if cmd.day != nil {
			d = int(cmd.day.IntValue())
		}
		if cmd.month != nil {
			m = int(cmd.month.IntValue())
		}
		if m == 0 {
			m = 1
		}

		choices = yearChoices(start, d, m)
	}

	cmd.ReplyAutocomplete(choices)
}

// executes when running the subcommand
func (cmd subcommandSet) interactionHandler() {

	authorID, err := strconv.ParseUint(cmd.user.ID, 10, 64)
	if err != nil {
		log.Printf("Error on parse author id of birthday command: %v\n", err)
		cmd.ReplyError()
		return
	}

	b := birthdayEntry{
		ID:      authorID,
		Day:     int(cmd.day.IntValue()),
		Month:   int(cmd.month.IntValue()),
		Visible: true,
	}
	if cmd.year != nil {
		b.Year = int(cmd.year.IntValue())
	}
	if cmd.visible != nil {
		b.Visible = cmd.visible.BoolValue()
	}

	hasBDay, err := cmd.hasBirthday(b.ID)
	if err != nil {
		log.Printf("Error on getting birthday data: %v\n", err)
		cmd.ReplyError()
		return
	}

	if hasBDay {
		cmd.updateBirthday(b)
	} else {
		cmd.setBirthday(b)
	}
}
