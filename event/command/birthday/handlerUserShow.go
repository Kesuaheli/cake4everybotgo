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

package birthday

import (
	"log"
	"strconv"
)

func (cmd UserShow) handler() {
	targetID, err := strconv.ParseUint(cmd.data.TargetID, 10, 64)
	if err != nil {
		log.Printf("Error on parse target id of birthday user show command: %v\n", err)
		cmd.ReplyError()
		return
	}

	target := cmd.data.Resolved.Members[cmd.data.TargetID]
	target.User = cmd.data.Resolved.Users[cmd.data.TargetID]

	log.Println("mem", target)
	log.Println("usr", target.User)

	hasBDay, err := cmd.hasBirthday(targetID)
	if err != nil {
		log.Printf("Error on show birthday: %v\n", err)
		cmd.ReplyError()
		return
	}

	var day, month, year int
	if hasBDay {
		// overwrite hasBDay with 'visible' field from database
		day, month, year, hasBDay, err = cmd.getBirthday(targetID)
		if err != nil {
			log.Printf("Error on show birthday: %v", err)
			cmd.ReplyError()
			return
		}
	}

	name := target.User.Username
	if target.Nick != "" {
		name = target.Nick
	}

	if !hasBDay {
		cmd.ReplyHiddenf("%s didn't enter their birthday nor set it visible.", name)
		return
	}

	cmd.Replyf("Birthday of %s is on %d.%d.%d", name, day, month, year)
}
