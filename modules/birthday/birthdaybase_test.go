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

import "testing"

func Test_birthdayEntry_DOW(t *testing.T) {
	tests := []struct {
		name string
		b    birthdayEntry
		want int
	}{
		{name: "bday-1", b: birthdayEntry{Day: 31, Month: 7, Year: 2000}, want: 0},
		{name: "bday", b: birthdayEntry{Day: 1, Month: 8, Year: 2000}, want: 1},
		{name: "bday+6", b: birthdayEntry{Day: 7, Month: 8, Year: 2000}, want: 0},

		{name: "random_0", b: birthdayEntry{Day: 19, Month: 9, Year: 1868}, want: 5},
		{name: "random_1", b: birthdayEntry{Day: 12, Month: 6, Year: 2000}, want: 0},
		{name: "random_2", b: birthdayEntry{Day: 29, Month: 8, Year: 2093}, want: 5},
		{name: "random_3", b: birthdayEntry{Day: 13, Month: 7, Year: 2202}, want: 1},
		{name: "random_4", b: birthdayEntry{Day: 17, Month: 9, Year: 2514}, want: 0},
		{name: "random_5", b: birthdayEntry{Day: 4, Month: 6, Year: 2638}, want: 0},
		{name: "random_6", b: birthdayEntry{Day: 21, Month: 5, Year: 2711}, want: 6},
		{name: "random_7", b: birthdayEntry{Day: 28, Month: 10, Year: 2788}, want: 4},
		{name: "random_8", b: birthdayEntry{Day: 17, Month: 4, Year: 2793}, want: 5},
		{name: "random_9", b: birthdayEntry{Day: 10, Month: 7, Year: 2948}, want: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := birthdayEntry{
				ID:      tt.b.ID,
				Day:     tt.b.Day,
				Month:   tt.b.Month,
				Year:    tt.b.Year,
				Visible: tt.b.Visible,
			}
			if got := b.DOW(); got != tt.want {
				t.Errorf("birthdayEntry.DOW() = %v, want %v", got, tt.want)
			}
		})
	}
}
