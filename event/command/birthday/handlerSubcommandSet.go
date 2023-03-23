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
		case "day":
			cmd.day = opt
		case "month":
			cmd.month = opt
		case "year":
			cmd.year = opt
		case "visible":
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

	if cmd.day.Focused {
		i, _ := strconv.Atoi(cmd.day.Value.(string))
		choices = append(choices, dayChoice(i))
	} else if cmd.month.Focused {
		// jan, mar, may, jul, aug, oct, dec
		choices = append(choices,
			monthChoice(1),
			monthChoice(3),
			monthChoice(5),
			monthChoice(7),
			monthChoice(8),
			monthChoice(10),
			monthChoice(12),
		)
		if cmd.day.IntValue() < 31 {
			// apr, jun, sep, nov
			choices = append(choices,
				monthChoice(4),
				monthChoice(6),
				monthChoice(9),
				monthChoice(11),
			)
		}
		if cmd.day.IntValue() == 29 && (cmd.year.IntValue()%4) == 0 {
			// feb (leap year)
			choices = append(choices,
				monthChoice(2),
			)
		}
		if cmd.day.IntValue() < 29 {
			// feb
			choices = append(choices,
				monthChoice(2),
			)
		}
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
		id:      authorID,
		day:     int(cmd.day.IntValue()),
		month:   int(cmd.month.IntValue()),
		visible: true,
	}
	if cmd.year != nil {
		b.year = int(cmd.year.IntValue())
	}
	if cmd.visible != nil {
		b.visible = cmd.visible.BoolValue()
	}

	hasBDay, err := cmd.hasBirthday(b.id)
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
