package compose

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

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

func loadComposeState(name string) (*Config, error) {
	_, maximum, err := getMinimumMaximumForState(name)
	if err != nil {
		return nil, err
	}

	if maximum == 0 {
		return nil, nil
	}

	file, err := os.ReadFile(fmt.Sprintf(".hcstate/%s-%d", name, maximum))
	if err != nil {
		return nil, err
	}

	config, err := decodeComposeConfig(string(file))
	if err != nil {
		return nil, err
	}

	return config, err
}

func storeComposeConfig(config *Config) error {
	encodedConfig, err := encodeComposeConfig(config)

	if err != nil {
		return err
	}

	switch config.State.Storage {
	case "kubernetes":
		fmt.Printf("To be implemented: K8s storage provider for states.")
		break
	default:
		if err := storeComposeConfigLocal(config.State.Name, &encodedConfig); err != nil {
			return err
		}
	}

	return nil
}

func storeComposeConfigLocal(name string, encodedConfig *string) error {
	minimum, maximum, err := getMinimumMaximumForState(name)
	if err != nil {
		return err
	}

	maximum = maximum + 1

	if err := os.WriteFile(fmt.Sprintf(".hcstate/%s-%d", name, maximum), []byte(*encodedConfig), 0644); err != nil {
		return err
	}

	if minimum > maximum-10 {
		return nil
	}

	if err := os.Remove(fmt.Sprintf(".hcstate/%s-%d", name, minimum)); err != nil {
		return err
	}

	return nil
}

func getMinimumMaximumForState(name string) (int, int, error) {
	if _, err := os.Stat(".hcstate"); os.IsNotExist(err) {
		if err := os.Mkdir(".hcstate", os.ModePerm); err != nil {
			return -1, -1, err
		}
	}

	files, err := os.ReadDir(".hcstate")
	if err != nil {
		return -1, -1, err
	}

	minimum, maximum := math.MaxInt, 0

	r, _ := regexp.Compile(fmt.Sprintf("^%s-(\\d+)$", name))

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		matches := r.FindStringSubmatch(file.Name())
		if len(matches) == 0 {
			continue
		}

		number, err := strconv.Atoi(matches[1])
		if err != nil {
			return -1, -1, err
		}

		if number > maximum {
			maximum = number
		}

		if number < minimum {
			minimum = number
		}
	}

	return minimum, maximum, nil
}
