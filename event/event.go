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
	"github.com/bwmarrin/discordgo"
)

// Register registers all events, like commands.
func Register(s *discordgo.Session, guildID string) error {
	err := registerCommands(s, guildID)
	if err != nil {
		return err
	}

	return nil
}

// AddListeners adds all event handlers to the given session s.
func AddListeners(s *discordgo.Session) {
	addCommandListeners(s)
	addVoiceStateListeners(s)

	addYouTubeListeners(s)
	addScheduledTriggers(s)
}
