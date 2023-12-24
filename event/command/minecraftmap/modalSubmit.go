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
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	label            *discordgo.TextInput
	id               *discordgo.TextInput
	world            *discordgo.Button
	world_the_nether *discordgo.Button
	world_the_end    *discordgo.Button
)

func (cmd Chat) ModalHandler() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		cmd.InteractionUtil = util.InteractionUtil{Session: s, Interaction: i}
		cmd.member = i.Member
		cmd.user = i.User
		if i.Member != nil {
			cmd.user = i.Member.User
		} else if i.User != nil {
			cmd.member = &discordgo.Member{User: i.User}
		}

		id := strings.Split(i.ModalSubmitData().CustomID, ".")[1]
		switch id {
		case "create_marker":
			cmd.handleModalCreateMarker()
		}
	}
}

func (cmd Chat) handleModalCreateMarker() {
	data := cmd.Interaction.ModalSubmitData()
	cmd.parseComponentData(data.Components)
	log.Printf("label: %s, id: %s", label.Value, id.Value)
	markerLock.RLock()
	markerBuilder[cmd.user.ID].Label = label.Value
	markerBuilder[cmd.user.ID].ID = id.Value
	markerLock.RUnlock()

	cmd.markerBuilder()

}

func (cmd Chat) parseComponentData(components []discordgo.MessageComponent) {
	for _, c := range components {
		switch c.Type() {
		case discordgo.ActionsRowComponent:
			cmd.parseComponentData(c.(*discordgo.ActionsRow).Components)
		case discordgo.TextInputComponent:
			input := c.(*discordgo.TextInput)
			switch input.CustomID {
			case "label":
				label = input
			case "id":
				id = input
			}
		case discordgo.ButtonComponent:
			button := c.(*discordgo.Button)
			switch button.CustomID {
			case "world":
				world = button
			case "world_the_nether":
				world_the_nether = button
			case "world_the_end":
				world_the_end = button
			}
		}
	}
}
