package compose

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func RunUp(config *Config) error {
	for name, url := range config.Repositories {
		if err := addHelmRepository(name, url); err != nil {
			return err
		}
	}

	// @TODO proper job queue system
	//var wg sync.WaitGroup

	for name, release := range config.Releases {
		//wg.Add(1)
		installHelmRelease(name, &release)
	}

	//wg.Wait()

	return nil
}

func RunDown(config *Config) error {
	// @TODO proper job queue system
	//var wg sync.WaitGroup

	for name, release := range config.Releases {
		//wg.Add(1)
		uninstallHelmRelease(name, &release)
	}

	//wg.Wait()

	return nil
}

func processRevision(config *Config) error {
	encodedConfig, err := encodeComposeConfig(config)

	if err != nil {
		return fmt.Errorf("Couldn't encode compose config: %s", err.Error())
	}

	if _, err := os.Stat(".hcstate"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(".hcstate", os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	err = ioutil.WriteFile(".hcstate/1", []byte(encodedConfig), 0644)
	if err != nil {
		return err
	}

	return nil
}
