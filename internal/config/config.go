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
package config

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver"
	"gopkg.in/yaml.v2"
)

func findComposeConfig() []string {
	var files []string
	filenames := []string{
		"helm-compose.yaml",
		"helm-compose.yml",
		"helmcompose.yaml",
		"helm-compose.yml",
		"helmcompose.yaml",
		"helmcompose.yml",
		"helmcompose",
		"compose.yaml",
		"compose.yml",
	}

	filepath.WalkDir(".", func(s string, d fs.DirEntry, e error) error {
		file := filepath.Base(s)

		for _, filename := range filenames {
			if file == filename {
				files = append(files, file)
			}
		}
		return nil
	})
	return files
}

func ParseConfigFile(filename string) (*Config, error) {
	var files []string
	if filename == "" {
		files = findComposeConfig()
	} else if _, err := os.Stat(filename); err != nil {
		return nil, fmt.Errorf("provided configuration file not found")
	} else {
		files = []string{filename}
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no configuration file found")
	}

	if len(files) > 1 {
		return nil, fmt.Errorf("expects only one configuration file but found multiple: %v", files)
	}

	file, err := os.ReadFile(files[0])
	if err != nil {
		return nil, err
	}

	return parseConfig(file)
}

func parseConfig(data []byte) (*Config, error) {
	config := Config{}
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if config.Version == "" {
		return nil, fmt.Errorf("missing apiVersion in config")
	}

	version, err := semver.NewVersion(config.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to parse apiVersion: %s", config.Version)
	}

	if semver.MustParse("1.0").GreaterThan(version) {
		return nil, fmt.Errorf("helm compose requires at least apiVersion 1.0 but got %s", config.Version)
	}

	return &config, nil
}
