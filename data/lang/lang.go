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

package lang

import (
	"log"

	"github.com/spf13/viper"
)

var langs = map[string]*viper.Viper{}

// Load loads (and reloads) the language files defined in the global
// config.
func Load() {
	log.Println("Loading languages...")

	// clear all existing languages
	for k := range langs {
		delete(langs, k)
	}

	// read language files defined in global config
	for _, langName := range viper.GetStringSlice("languages") {
		if langName == "" {
			continue
		}
		// skip duplicates
		if _, ok := langs[langName]; ok {
			continue
		}

		log.Printf("Loading language '%s'...\n", langName)

		lang := viper.New()
		lang.SetConfigType("yaml")
		lang.AddConfigPath("data/lang")
		lang.SetConfigName(langName)

		err := lang.ReadInConfig()
		if err != nil {
			log.Printf("WARNING: Could not load language '%s': %v", langName, err)
			continue
		}
		langs[langName] = lang
	}

	if len(langs) == 0 {
		log.Fatalln("Could not load languages: needs at least one loaded language!")
	}

	log.Printf("Loaded %d language(s), with %s beeing the fallback language!\n", len(langs), FallbackLang())
}

// Get returns the configured translation for key in the given
// language lang.
//   - If lang is not a loaded language, Get translates key with the
//     fallback language.
//   - If key either does not exits or is an empty string, Get
//     translates key with the fallback language.
//   - However, if lang already is the fallback language in one of the
//     cases above, Get returns key instead.
//
// In all three of these 'fail cases', Get will print a warning
// message in the log
func Get(key, lang string) string {
	if len(langs) == 0 {
		log.Println()
		log.Printf("ERROR: Tried to get translation, but no language loaded")
		log.Println()
		return key
	}

	v, ok := langs[lang]
	fLang := FallbackLang()
	if !ok {
		if lang == fLang {
			log.Println()
			log.Printf("ERROR: Tried to get key from fallback language ('%s'), but its not load", fLang)
			log.Println()
			return key
		}
		log.Printf("WARNING: language '%s' is not loaded, using '%s' as fallback instead", lang, fLang)
		return Get(key, fLang)
	}

	val := v.GetString(key)
	if val != "" {
		return val
	}

	if lang != fLang {
		log.Printf("WARNING: key '%s' is not defined in fallback language '%s'", key, lang)
		return key
	}
	log.Printf("WARNING: key '%s' is not defined in language '%s', using '%s' as fallback instead", key, lang, fLang)
	return Get(key, fLang)
}

// FallbackLang returns the bots fallback language, which is the
// first loaded language string in the 'languages' list from the
// global config.
//
// When 'languages' is empty or only has empty string entries, it
// fails and calls os.Exit(1).
func FallbackLang() string {
	for _, langName := range viper.GetStringSlice("languages") {
		if langName == "" {
			continue
		}
		if _, ok := langs[langName]; !ok {
			continue
		}
		return langName
	}

	log.Fatalln("ERROR: No languages in config")
	return ""
}
