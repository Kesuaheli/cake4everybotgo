// Copyright 2022-2023 Kesuaheli
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

package event

import (
	"cake4everybot/event/command"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func registerCommands(s *discordgo.Session, guildID string) error {

	var commandsList []command.Command
	// This is the list of commands to use. Add a command via simply
	// appending the struct (which must implement the interface
	// command.Command) to the list, i.e.:
	// commandsList = append(commandsList, command.MyCommand{})
	commandsList = append(commandsList, command.Birthday{})

	// early return when there're no commands to add, and remove all previously registered commands
	if len(commandsList) == 0 {
		removeUnusedCommands(s, guildID, nil)
		return nil
	}

	// make an array of ApplicationCommands and perform a bulk change using it
	appCommandsList := make([]*discordgo.ApplicationCommand, 0, len(commandsList))
	for _, cmd := range commandsList {
		appCommandsList = append(appCommandsList, cmd.AppCmd())
		command.CommandMap[cmd.AppCmd().Name] = cmd
	}
	commandNames := make([]string, 0, len(command.CommandMap))
	for k := range command.CommandMap {
		commandNames = append(commandNames, k)
	}
	fmt.Printf("Adding used commands: %v...\n", commandNames)
	createdCommands, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, guildID, appCommandsList)
	if err != nil {
		return fmt.Errorf("failed on bulk overwrite commands: %v", err)
	}

	removeUnusedCommands(s, guildID, createdCommands)
	return err
}

func removeUnusedCommands(s *discordgo.Session, guildID string, createdCommands []*discordgo.ApplicationCommand) {
	allRegisteredCommands, err := s.ApplicationCommands(s.State.User.ID, guildID)
	if err != nil {
		fmt.Printf("Error while removing unused commands: Could not get registered commands from guild '%s'. Err: %v\n", guildID, err)
		return
	}
	newCmdIds := make(map[string]bool)
	for _, cmd := range createdCommands {
		newCmdIds[cmd.ID] = true
	}

	// Find unused commands by iterating over all commands and check if the ID is in the currently registered commands. If not, remove it.
	for _, cmd := range allRegisteredCommands {
		if !newCmdIds[cmd.ID] {
			err = s.ApplicationCommandDelete(s.State.User.ID, guildID, cmd.ID)
			if err != nil {
				fmt.Printf("Error while removing unused commands: Could not delete comand '%s' ('/%s'). Err: %v\n", cmd.ID, cmd.Name, err)
				continue
			}
			fmt.Printf("Removed unused command: '/%s'\n", cmd.Name)
		}
	}

}

func addCommandListeners(s *discordgo.Session) {
	s.AddHandler(func(s *discordgo.Session, event *discordgo.InteractionCreate) {
		if cmd, ok := command.CommandMap[event.ApplicationCommandData().Name]; ok {
			cmd.CmdHandler()(s, event)
		}
	})

}
