package compose

import "sync"

func RunUp(config *Config) error {
	for name, url := range config.Repositories {
		if err := addHelmRepository(name, url); err != nil {
			return err
		}
	}

	// @TODO proper job queue system
	var wg sync.WaitGroup

	for name, release := range config.Releases {
		wg.Add(1)
		go installHelmRelease(name, &release, &wg)
	}

	wg.Wait()

	return nil
}
