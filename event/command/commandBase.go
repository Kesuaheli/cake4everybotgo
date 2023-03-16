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

type Command interface {
	// Definition of a command.
	// E.g., name, description, options, subcommands.
	AppCmd() *discordgo.ApplicationCommand

	// Function of a command.
	// All things that should happen at execution.
	CmdHandler() func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var CommandMap map[string]Command

func init() {
	CommandMap = make(map[string]Command)
}
