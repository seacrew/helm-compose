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
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver"
	"gopkg.in/yaml.v2"
)

var (
	V1_0 = semver.MustParse("1.0")
	V1_1 = semver.MustParse("1.1")
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

func ParseComposeFile(filename string) (*Config, error) {
	var files []string
	if filename == "" {
		files = findComposeConfig()
	} else if filename == "-" {
		reader := bufio.NewReader(os.Stdin)

		data := []byte{}
		for {
			b, err := reader.ReadBytes('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			data = append(data, b...)
		}

		return parseComposeData(data)
	} else if _, err := os.Stat(filename); err != nil {
		return nil, fmt.Errorf("provided compose file not found")
	} else {
		files = []string{filename}
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no compose file found")
	}

	if len(files) > 1 {
		return nil, fmt.Errorf("expected only one compose file but found multiple: %v", files)
	}

	file, err := os.ReadFile(files[0])
	if err != nil {
		return nil, err
	}

	return parseComposeData(file)
}

func parseComposeData(data []byte) (*Config, error) {
	config := Config{}
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if err := validateCompose(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateCompose(config *Config) error {
	if config.Version == "" {
		return fmt.Errorf("missing apiVersion in config")
	}

	version, err := semver.NewVersion(config.Version)
	if err != nil {
		return fmt.Errorf("failed to parse apiVersion: %s", config.Version)
	}

	if version.LessThan(V1_0) {
		return fmt.Errorf("helm compose requires at least apiVersion 1.0 but got %s", config.Version)
	}

	if err := validateComposeFeatures(version, config); err != nil {
		return err
	}

	return nil
}

func validateComposeFeatures(version *semver.Version, config *Config) error {
	if err := validateCompose1_1(version, config); err != nil {
		return fmt.Errorf("apiVersion 1.1+ necessary: %s", err)
	}

	return nil
}

func validateCompose1_1(version *semver.Version, config *Config) error {
	if version.GreaterThan(V1_0) {
		return nil
	}

	for name, release := range config.Releases {
		if release.Wait {
			return fmt.Errorf("trying to use 'wait' in release '%s'", name)
		}
	}

	return nil
}
