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
	"github.com/bwmarrin/discordgo"

	"cake4everybot/event/command/util"
)

// A user command of the birthday package. It
// adds the ability to directly show a users
// birthday through a simple context click.
type UserShow struct {
	birthdayBase

	data discordgo.ApplicationCommandInteractionData
}

func (cmd UserShow) AppCmd() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Type: discordgo.UserApplicationCommand,
		Name: "show birthday",
	}
}

func (cmd UserShow) CmdHandler() func(s *discordgo.Session, i *discordgo.InteractionCreate) {

	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		cmd.InteractionUtil = util.InteractionUtil{Session: s, Interaction: i}
		cmd.member = i.Member
		cmd.user = i.User
		if i.Member != nil {
			cmd.user = i.Member.User
		}

		cmd.data = cmd.Interaction.ApplicationCommandData()
		cmd.handler()
	}
}
