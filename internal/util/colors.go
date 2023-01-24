package util

import (
	"fmt"
	"hash/fnv"
	"strconv"

	"github.com/jwalton/go-supportscolor"
)

var colors [](func(...interface{}) string) = [](func(...interface{}) string){
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

func HashColor(s string) func(...interface{}) string {
	if !supportscolor.Stdout().SupportsColor {
		return colors[0]
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

		if subtotal < len(colors) {
			return colors[subtotal]
		}

		hash = fmt.Sprint(subtotal)
	}
}
