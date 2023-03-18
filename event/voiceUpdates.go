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

package event

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var NO_MIC_CHANNEL_ID string

func addVoiceStateListeners(s *discordgo.Session) {
	handler := func(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {

		isAfk, err := isAfkVoiceChannel(s, e)
		if err != nil {
			fmt.Print(err)
			return
		}

		if isVoiceChannelUpdate(e) && (e.ChannelID == "" || isAfk) {
			setNoMicPermission(s, e, false)
		} else if isVoiceChannelUpdate(e) {
			setNoMicPermission(s, e, true)
		}
	}

	s.AddHandler(handler)
}

func isVoiceChannelUpdate(e *discordgo.VoiceStateUpdate) bool {
	return e.BeforeUpdate == nil || e.BeforeUpdate.ChannelID != e.ChannelID
}

func isAfkVoiceChannel(s *discordgo.Session, e *discordgo.VoiceStateUpdate) (bool, error) {
	guild, err := s.Guild(e.GuildID)
	if err != nil {
		return false, fmt.Errorf("ERROR: on join afk vc: %v\n", err)
	}

	return e.ChannelID == guild.AfkChannelID, nil
}

func setNoMicPermission(s *discordgo.Session, e *discordgo.VoiceStateUpdate, state bool) {
	var err error
	if state {
		err = s.ChannelPermissionSet(NO_MIC_CHANNEL_ID,
			e.Member.User.ID,
			discordgo.PermissionOverwriteTypeMember,
			discordgo.PermissionViewChannel,
			0)
	} else {
		err = s.ChannelPermissionDelete(NO_MIC_CHANNEL_ID, e.Member.User.ID)
	}

	if err != nil {
		fmt.Printf("Error on no mic permission: %v\n", err)
		return
	}
}
