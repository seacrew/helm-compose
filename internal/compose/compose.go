package compose

import (
	"github.com/nileger/helm-compose/internal/util"
)

func RunUp(config *Config) error {

	for name, url := range config.Repositories {
		if err := util.AddHelmRepository(name, url); err != nil {
			return err
		}
	}

	return nil
}
