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
	"math"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	re = regexp.MustCompile(`\$\{(.*?)\}`)
)

func IsDebug() bool {
	return os.Getenv("HELM_DEBUG") == "true"
}

func DebugPrint(format string, a ...interface{}) {
	if IsDebug() {
		fmt.Printf(format+"\n", a...)
	}
}

func Execute(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()

	if output == nil {
		return "", err
	}

	text := string(output)

	if len(text) > 5 && text[:6] == "Error:" {
		text = strings.TrimSpace(text[6:])
	}

	return text, err
}

func ConvertJson(obj interface{}) interface{} {
	switch c := obj.(type) {
	case map[interface{}]interface{}:
		m := map[string]interface{}{}
		for k, v := range c {
			m[k.(string)] = ConvertJson(v)
		}
		return m
	case []interface{}:
		for k, v := range c {
			c[k] = ConvertJson(v)
		}
	case string:
		str := obj.(string)
		matches := re.FindStringSubmatch(str)
		for _, match := range matches {
			str = strings.Replace(str, "${"+match+"}", os.Getenv(match), 1)
		}

		return str
	}
	return obj
}

func MinMax(ints []int) (int, int) {
	if len(ints) == 0 {
		return 0, 0
	}

	minimum, maximum := math.MaxInt, 0

	for _, i := range ints {
		if i > maximum {
			maximum = i
		}

		if i < minimum {
			minimum = i
		}
	}

	return minimum, maximum
}
