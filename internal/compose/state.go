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
package compose

import (
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
)

func loadComposeState(config *Config) (*Config, error) {
	name := config.Name
	if len(config.Storage.NameOverride) > 0 {
		name = config.Storage.NameOverride
	}

	var err error
	var data []byte

	switch config.Storage.Type {
	case "local", "":
		data, err = loadComposeStateLocal(name, config.Storage.Path)
		if err != nil {
			return nil, err
		}
	case "kubernetes":
		fmt.Printf("To be implemented: K8s storage provider for states.")
		break
	default:
		return nil, errors.New(fmt.Sprintf("Storage type of '%s' does not exists.", config.Storage.Type))
	}

	if data == nil {
		return nil, nil
	}

	prevConfig, err := decodeComposeConfig(string(data))
	if err != nil {
		return nil, err
	}

	return prevConfig, nil
}

func loadComposeStateLocal(name string, path string) ([]byte, error) {
	if len(path) == 0 {
		path = ".hcstate"
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			return nil, err
		}
	}

	_, maximum, err := getMinMaxStateLocal(name, path)
	if err != nil {
		return nil, err
	}

	if maximum == 0 {
		return nil, nil
	}

	file, err := os.ReadFile(fmt.Sprintf("%s/%s-%d", path, name, maximum))
	if err != nil {
		return nil, err
	}

	return file, nil
}

func storeComposeConfig(config *Config) error {
	encodedConfig, err := encodeComposeConfig(config)

	name := config.Name
	if len(config.Storage.NameOverride) > 0 {
		name = config.Storage.NameOverride
	}

	if config.Storage.NumberOfStates <= 0 {
		config.Storage.NumberOfStates = 10
	}

	if err != nil {
		return err
	}

	switch config.Storage.Type {
	case "local", "":
		if err := storeComposeConfigLocal(name, config, &encodedConfig); err != nil {
			return err
		}
	case "kubernetes":
		fmt.Printf("To be implemented: K8s storage provider for states.")
		break
	default:
		return errors.New(fmt.Sprintf("Storage type of '%s' does not exists.", config.Storage.Type))
	}

	return nil
}

func storeComposeConfigLocal(name string, config *Config, encodedConfig *string) error {
	path := config.Storage.Path
	if len(path) == 0 {
		path = ".hcstate"
	}

	minimum, maximum, err := getMinMaxStateLocal(name, path)
	if err != nil {
		return err
	}

	maximum = maximum + 1

	if err := os.WriteFile(fmt.Sprintf("%s/%s-%d", path, name, maximum), []byte(*encodedConfig), 0644); err != nil {
		return err
	}

	if minimum > maximum-config.Storage.NumberOfStates {
		return nil
	}

	if err := os.Remove(fmt.Sprintf("%s/%s-%d", path, name, minimum)); err != nil {
		return err
	}

	return nil
}

func getMinMaxState(states []int) (int, int) {
	if len(states) == 0 {
		return 0, 0
	}

	minimum, maximum := math.MaxInt, 0

	for _, state := range states {
		if state > maximum {
			maximum = state
		}

		if state < minimum {
			minimum = state
		}
	}

	return minimum, maximum
}

func getMinMaxStateLocal(name string, path string) (int, int, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return -1, -1, err
	}

	r, _ := regexp.Compile(fmt.Sprintf("^%s-(\\d+)$", name))

	states := []int{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		matches := r.FindStringSubmatch(file.Name())
		if len(matches) == 0 {
			continue
		}

		state, err := strconv.Atoi(matches[1])
		if err != nil {
			return -1, -1, nil
		}

		states = append(states, state)
	}

	minimum, maximum := getMinMaxState(states)
	return minimum, maximum, nil
}
