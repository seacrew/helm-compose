/*
Copyright © 2023 The Helm Compose Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package util

import (
	"fmt"
	"hash/fnv"
	"strconv"

	"github.com/jwalton/go-supportscolor"
)

type ColorPrinter struct {
	colorFunc func(...interface{}) string
}

func NewColorPrinter(s string) *ColorPrinter {
	c := hashColor(s)
	return &ColorPrinter{
		colorFunc: c,
	}
}

func (c ColorPrinter) Printf(format string, a ...any) {
	fmt.Printf(c.colorFunc(format)+"\n", a...)
}

var colorFuncs [](func(...interface{}) string) = [](func(...interface{}) string){
	color("%s"), // fallback
	color("\033[1;32m%s\033[0m"),
	color("\033[1;33m%s\033[0m"),
	color("\033[1;34m%s\033[0m"),
	color("\033[1;35m%s\033[0m"),
	color("\033[1;36m%s\033[0m"),
	color("\033[1;90m%s\033[0m"),
	color("\033[1;92m%s\033[0m"),
	color("\033[1;93m%s\033[0m"),
	color("\033[1;94m%s\033[0m"),
	color("\033[1;95m%s\033[0m"),
	color("\033[1;96m%s\033[0m"),
}

func color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

func hashColor(s string) func(...interface{}) string {
	if !supportscolor.Stdout().SupportsColor {
		return colorFuncs[0]
	}

	h := fnv.New32a()
	h.Write([]byte(s))

	hash := fmt.Sprint(h.Sum32())
	for {
		subtotal := 0
		for _, r := range hash {
			value, _ := strconv.Atoi(string(r))
			subtotal += value
		}

		if subtotal < len(colorFuncs) {
			return colorFuncs[subtotal]
		}

		hash = fmt.Sprint(subtotal)
	}
}
