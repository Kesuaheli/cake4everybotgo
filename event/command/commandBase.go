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

package command

import "github.com/bwmarrin/discordgo"

// Command is an interface wrapper for all commands. Including chat-
// comamnds (slash-commands), message-commands, and user-commands.
type Command interface {
	// Definition of a command.
	// E.g., name, description, options, subcommands.
	AppCmd() *discordgo.ApplicationCommand

	// Function of a command.
	// All things that should happen at execution.
	CmdHandler() func(s *discordgo.Session, i *discordgo.InteractionCreate)

	// Sets the registered command ID for internal uses after uploading to discord
	SetID(id string)
	// Gets the registered command ID
	GetID() string
}

// CommandMap holds all active commands. It maps them from a unique
// name identifier to the corresponding Command.
//
// Here the name is used, because Discord uses the name too to
// identify seperate commands. When a command is beeing registered
// with a name that already is beeing registerd as a command by this
// application (bot), then the new one will simply overwrite it and
// automatically ungerister the old one.
var CommandMap map[string]Command

func init() {
	CommandMap = make(map[string]Command)
}
