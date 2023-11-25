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

package util

// ContainsInt reports whether at least one of num is at least once
// anywhere in i.
func ContainsInt(i []int, num ...int) bool {
	for _, x := range i {
		for _, y := range num {
			if x == y {
				return true
			}
		}
	}
	return false
}

// ContainsString reports whether at least one of str is at least
// once anywhere in s.
func ContainsString(s []string, str ...string) bool {
	for _, x := range s {
		for _, y := range str {
			if x == y {
				return true
			}
		}
	}
	return false
}

// Btoi returns the integer for the given boolean b.
//
//	Btoi(false) => 0
//	Btoi(true) => 1
func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
