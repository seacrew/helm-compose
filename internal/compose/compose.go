package compose

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
