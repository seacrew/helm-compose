package compose

import (
	"fmt"
	"io/fs"
	"io/ioutil"
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

func ParseConfig(filename string) (*Config, error) {
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

	file, err := ioutil.ReadFile(files[0])

	if err != nil {
		return nil, err
	}

	config := Config{}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	if config.Version == "" {
		return nil, fmt.Errorf("missing apiVersion in file %s", files[0])
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
