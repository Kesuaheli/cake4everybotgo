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

package twitch

import (
	logger "log"

	"github.com/kesuaheli/twitchgo"
)

var log logger.Logger = *logger.New(logger.Writer(), "[Twitch] ", logger.LstdFlags|logger.Lmsgprefix)

// MessageHandler handles new messages from the twitch chat(s). It will be called on every new
// message.
func MessageHandler(t *twitchgo.Twitch, channel string, user *twitchgo.User, message string) {
	log.Printf("<%s@%s> %s", user.Nickname, channel, message)
}
