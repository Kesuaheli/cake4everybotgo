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
	"time"

	"cake4everybot/event/command/birthday"
	webYT "cake4everybot/webserver/youtube"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

func addScheduledTriggers(s *discordgo.Session, webChan chan struct{}) {
	go scheduleBirthdayCheck(s)
	go refreshYoutube(webChan)
}

func scheduleBirthdayCheck(s *discordgo.Session) {
	HOUR := viper.GetInt("event.birthday_hour")

	time.Sleep(time.Second * 5)
	for {
		now := time.Now()

		nextRun := time.Date(now.Year(), now.Month(), now.Day(), HOUR, 0, 0, 0, now.Location())
		if nextRun.Before(now) {
			nextRun = nextRun.Add(time.Hour * 24)
		}
		time.Sleep(nextRun.Sub(now))

		birthday.Check(s)
	}
}

func refreshYoutube(webChan chan struct{}) {
	<-webChan
	for {
		webYT.RefreshSubscriptions()

		// loop every 4 days
		time.Sleep(4 * 24 * time.Hour)
	}
}
