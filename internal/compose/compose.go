package compose

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

	// @TODO proper job queue system
	//var wg sync.WaitGroup

	for name, release := range config.Releases {
		//wg.Add(1)
		installHelmRelease(name, &release)
	}

	if previousConfig == nil {
		return nil
	}

	for name, release := range previousConfig.Releases {
		if _, ok := config.Releases[name]; ok {
			continue
		}

		uninstallHelmRelease(name, &release)
	}

	//wg.Wait()

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

	// @TODO proper job queue system
	//var wg sync.WaitGroup

	for name, release := range config.Releases {
		//wg.Add(1)
		uninstallHelmRelease(name, &release)
	}

	//wg.Wait()

	return nil
}
