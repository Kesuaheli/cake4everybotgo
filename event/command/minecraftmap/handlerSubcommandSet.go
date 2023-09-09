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

package minecraftmap

import (
	"cake4everybot/data/lang"

	"github.com/bwmarrin/discordgo"
)

// The set subcommand. Used when executing the
// slash-command "/birthday set".
type subcommandSet struct {
	Chat
	*discordgo.ApplicationCommandInteractionDataOption

	set   *discordgo.ApplicationCommandInteractionDataOption // required
	id    *discordgo.ApplicationCommandInteractionDataOption // required
	label *discordgo.ApplicationCommandInteractionDataOption // required
	world *discordgo.ApplicationCommandInteractionDataOption // required
	posX  *discordgo.ApplicationCommandInteractionDataOption // required
	posY  *discordgo.ApplicationCommandInteractionDataOption // required
	posZ  *discordgo.ApplicationCommandInteractionDataOption // required
	icon  *discordgo.ApplicationCommandInteractionDataOption // required
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
		case lang.GetDefault(tp + "option.set.option.set"):
			cmd.set = opt
		case lang.GetDefault(tp + "option.set.option.id"):
			cmd.id = opt
		case lang.GetDefault(tp + "option.set.option.label"):
			cmd.label = opt
		case lang.GetDefault(tp + "option.set.option.world"):
			cmd.world = opt
		case lang.GetDefault(tp + "option.set.option.posX"):
			cmd.posX = opt
		case lang.GetDefault(tp + "option.set.option.posY"):
			cmd.posY = opt
		case lang.GetDefault(tp + "option.set.option.posZ"):
			cmd.posZ = opt
		case lang.GetDefault(tp + "option.set.option.icon"):
			cmd.icon = opt
		}
	}

	switch cmd.Interaction.Type {
	case discordgo.InteractionApplicationCommand:
		cmd.interactionHandler()
	}
}

// executes when running the subcommand
func (cmd subcommandSet) interactionHandler() {
	cmd.ReplyHidden("W.I.P.")
}
