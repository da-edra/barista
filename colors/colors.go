// Copyright 2017 Google Inc.
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

// Package colors provides helper functions to manage color and color schemes.
package colors

import (
	"bufio"
	"image/color"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/spf13/afero"
)

// Hex sanity-checks and constructs a color from a hex-string.
// Any string that can be parsed by colorful is acceptable.
func Hex(hex string) color.Color {
	c, err := colorful.Hex(hex)
	if err != nil {
		return nil
	}
	return c
}

// Scheme gets a color from the user-defined color scheme.
// Some common names are 'good', 'bad', and 'degraded'.
func Scheme(name string) color.Color {
	return scheme[name]
}

// scheme holds the mapping of "name" to colour values.
// Modules can use this to provide default colors for their output
// by using the commonly accepted names "good", "bad", and "degraded".
// Bar authors can also define arbitrary names, e.g. to load XResource based colours
// from i3 using the "LoadFromArgs" method.
var scheme = map[string]color.Color{}

func splitAtLastEqual(s string) (string, string, bool) {
	idx := strings.LastIndex(s, "=")
	if idx < 0 {
		return "", "", false
	}
	return s[:idx], s[idx+1:], true
}

// LoadFromArgs loads a color scheme from command-line arguments of the form name=value.
func LoadFromArgs(args []string) {
	for _, arg := range args {
		if name, value, ok := splitAtLastEqual(arg); ok {
			if color := Hex(value); color != nil {
				scheme[name] = color
			}
		}
	}
}

// LoadFromMap sets the colour scheme from code.
func LoadFromMap(s map[string]string) {
	for name, value := range s {
		if color := Hex(value); color != nil {
			scheme[name] = color
		}
	}
}

var fs = afero.NewOsFs()

// LoadFromConfig loads a color scheme from a i3status config file
// by reading all keys that start with 'color_'
func LoadFromConfig(filename string) error {
	f, err := fs.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if !strings.HasPrefix(line, "color_") {
			continue
		}
		name, value, ok := splitAtLastEqual(line)
		if !ok {
			continue
		}
		name = strings.TrimSpace(name[len("color_"):])
		value = strings.TrimSpace(value)
		if value[0] == '"' && value[len(value)-1] == '"' {
			value = value[1 : len(value)-1]
		} else if value[0] == '\'' && value[len(value)-1] == '\'' {
			value = value[1 : len(value)-1]
		}
		if color := Hex(value); color != nil {
			scheme[name] = color
		}
	}
	return nil
}
