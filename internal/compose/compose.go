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
