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
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"cake4everybot/event/command/util"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

const (
	// Prefix for translation key, i.e.:
	//   key := tp+"base" // => minecraftmap
	tp = "discord.command.minecraftmap."
)

type mapBase struct {
	util.InteractionUtil
	member *discordgo.Member
	user   *discordgo.User
}

type marker struct {
	Set    string
	ID     string `json:"id"`
	Label  string `json:"label"`
	World  string `json:"world"`
	PosX   int    `json:"posX"`
	PosY   int    `json:"posY"`
	PosZ   int    `json:"posZ"`
	IconID string `json:"icon_id"`
}

// Returns a readable Form of the marker
func (m marker) String() string {
	return fmt.Sprintf("%s [id: %s/%s], %d, %d, %d, icon: %s", m.Label, m.Set, m.ID, m.PosX, m.PosY, m.PosZ, m.IconID)
}

func (m marker) post() error {
	url := viper.GetString("minecraft.map.url") + "/marker/" + m.Set

	buf, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("marker Marshal Error: %s", err)
	}
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(buf))
	if err != nil {
		return fmt.Errorf("marker Post Error: %s", err)
	}

	if r.Response.StatusCode != 200 {
		return fmt.Errorf("got wrong response: %d: %v+", r.Response.StatusCode, r.Response.Body)
	}

	return nil
}

type set struct {
	ID      string   `json:"id"`
	Label   string   `json:"label"`
	Markers []string `json:"markers,omitempty"`
}

// Returns a readable Form of the set
func (s set) String() string {
	return fmt.Sprintf("%s [id: %s], icon: %v", s.Label, s.ID, s.Markers)
}
