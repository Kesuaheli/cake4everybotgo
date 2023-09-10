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
	"github.com/bwmarrin/discordgo"
)

// The marker subcommand. Used when executing the
// slash-command "/minecraftmap marker".
type subcommandMarker struct {
	Chat
	*discordgo.ApplicationCommandInteractionDataOption
}

// Constructor for subcommandmarker, the struct for
// the slash-command "/minecraftmap marker".
func (cmd Chat) subcommandMarker() subcommandMarker {
	subcommand := cmd.Interaction.ApplicationCommandData().Options[0]
	return subcommandMarker{
		Chat:                                    cmd,
		ApplicationCommandInteractionDataOption: subcommand,
	}
}

func (cmd subcommandMarker) handler() {
	switch cmd.Interaction.Type {
	case discordgo.InteractionApplicationCommand:
		cmd.interactionHandler()
	}
}

// executes when running the subcommand
func (cmd subcommandMarker) interactionHandler() {
	cmd.ReplyModal(tp, "create_marker", cmd.create_marker_id()...)
	cmd.ReplyHidden("W.I.P.")
}
