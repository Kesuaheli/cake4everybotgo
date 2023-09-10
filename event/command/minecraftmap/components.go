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
	"cake4everybot/event/command/util"

	"github.com/bwmarrin/discordgo"
)

func (Chat) create_marker_id() []discordgo.MessageComponent {
	component_group_id := "create_marker"
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			util.CreateTextInputComponent(tp, component_group_id, "label", discordgo.TextInputShort, true, 4, 25),
		}},
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			util.CreateTextInputComponent(tp, component_group_id, "id", discordgo.TextInputShort, true, 4, 25),
		}},
	}
}

func (Chat) create_marker_world() []discordgo.MessageComponent {
	component_group_id := "create_marker"
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			util.CreateButtonComponent(tp, component_group_id, "world", discordgo.SuccessButton, discordgo.ComponentEmoji{Name: "🗺️" /*:map:*/}),
			util.CreateButtonComponent(tp, component_group_id, "world_the_nether", discordgo.DangerButton, discordgo.ComponentEmoji{Name: "🔥" /*:fire:*/}),
			util.CreateButtonComponent(tp, component_group_id, "world_the_end", discordgo.SecondaryButton, discordgo.ComponentEmoji{Name: "🌌" /*:milky_way:*/}),
		}},
	}
}

func (Chat) create_marker_position() []discordgo.MessageComponent {
	component_group_id := "create_marker"
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			util.CreateTextInputComponent(tp, component_group_id, "posx", discordgo.TextInputShort, true, 0, 7),
			util.CreateTextInputComponent(tp, component_group_id, "posy", discordgo.TextInputShort, true, 0, 7),
			util.CreateTextInputComponent(tp, component_group_id, "posz", discordgo.TextInputShort, true, 0, 7),
		}},
	}
}

func (Chat) create_marker_icon() []discordgo.MessageComponent {
	component_group_id := "create_marker"
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			util.CreateTextInputComponent(tp, component_group_id, "icon", discordgo.TextInputShort, true, 0, 20),
		}},
	}

}
