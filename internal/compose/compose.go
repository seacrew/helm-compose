package compose

import "sync"

func RunUp(config *Config) error {
	for name, url := range config.Repositories {
		if err := addHelmRepository(name, url); err != nil {
			return err
		}
	}

	previousConfig, err := loadComposeState(config.State.Name)
	if err != nil {
		return err
	}

	if err := storeComposeConfig(config); err != nil {
		return err
	}

	var wg sync.WaitGroup

	for name, release := range config.Releases {
		wg.Add(1)
		go func(name string, release Release) {
			installHelmRelease(name, &release)
			wg.Done()
		}(name, release)
	}

	if previousConfig == nil {
		wg.Wait()
		return nil
	}

	for name, release := range previousConfig.Releases {
		wg.Add(1)
		go func(name string, release Release) {
			if _, ok := config.Releases[name]; ok {
				wg.Done()
				return
			}

			uninstallHelmRelease(name, &release)
			wg.Done()
		}(name, release)
	}

	wg.Wait()

	return nil
}

func RunDown(config *Config) error {
	previousConfig, err := loadComposeState(config.State.Name)
	if err != nil {
		return err
	}

	if previousConfig != nil {
		config = previousConfig
	}

	var wg sync.WaitGroup

	for name, release := range config.Releases {
		wg.Add(1)
		go func(name string, release Release) {
			if _, ok := config.Releases[name]; ok {
				wg.Done()
				return
			}

			uninstallHelmRelease(name, &release)
			wg.Done()
		}(name, release)
	}

	wg.Wait()

	return nil
}
