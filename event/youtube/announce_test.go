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

package youtube

import "testing"

func Test_saveTrimString(t *testing.T) {
	type args struct {
		s string
		n int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"fewer then n chars", args{"cool", 10}, "cool"},
		{"fewer then n chars", args{"this is pretty awesome", 99}, "this is pretty awesome"},
		{"emtpy", args{"", 5}, ""},
		{"emtpy", args{" ", 5}, ""},
		{"one long word", args{"Llanfairpwllgwyngyllgogerychwyrndrobwllllantysiliogogogoch", 10}, "Llanfairpwllgwyngyllgogerychwyrndrobwllllantysiliogogogoch"},
		{"cut words", args{"this is pretty awesome", 10}, "this is ..."},
		{"at spot", args{"this is pretty awesome", 22}, "this is pretty awesome"},
		{"just below at spot", args{"this is pretty awesome", 19}, "this is pretty awesome"},
		{"prevent longer output edge case", args{"ThisIsALongWordWithASingleLetterAfterASpacCharacter s", 20}, "ThisIsALongWordWithASingleLetterAfterASpacCharacter s"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := saveTrimText(tt.args.s, tt.args.n); got != tt.want {
				t.Errorf("saveTrimString() = %v, want %v", got, tt.want)
			}
		})
	}
}
