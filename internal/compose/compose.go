package compose

func RunUp(config *Config) error {
	for name, url := range config.Repositories {
		if err := addHelmRepository(name, url); err != nil {
			return err
		}
	}

	for name, release := range config.Releases {
		err := installHelmRelease(name, &release)
		if err != nil {
			return err
		}
	}

	return nil
}
