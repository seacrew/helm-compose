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

import (
	"fmt"
	"sync"

	cfg "github.com/seacrew/helm-compose/internal/config"
	"github.com/seacrew/helm-compose/internal/state"
)

func RunUp(config *cfg.Config) error {
	for name, url := range config.Repositories {
		if err := addHelmRepository(name, url); err != nil {
			return err
		}
	}

	previousConfig, err := state.Load(config)
	if err != nil {
		return err
	}

	if err := state.Store(config); err != nil {
		return err
	}

	var wg sync.WaitGroup

	for name, release := range config.Releases {
		wg.Add(1)
		go func(name string, release cfg.Release) {
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
		go func(name string, release cfg.Release) {
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

func RunDown(config *cfg.Config) error {
	previousConfig, err := state.Load(config)
	if err != nil {
		return err
	}

	if previousConfig != nil {
		config = previousConfig
	}

	var wg sync.WaitGroup

	for name, release := range config.Releases {
		wg.Add(1)
		go func(name string, release cfg.Release) {
			uninstallHelmRelease(name, &release)
			wg.Done()
		}(name, release)
	}

	wg.Wait()

	return nil
}

func ListRevisions(config *cfg.Config) error {
	revisions, err := state.List(config)
	if err != nil {
		return err
	}

	fmt.Printf("| Date             | Revision |\n")
	fmt.Printf("| ---------------- | -------- |\n")
	for _, rev := range revisions {
		fmt.Printf("| %d-%02d-%02d %02d:%02d | %8d |\n",
			rev.DateTime.Year(), rev.DateTime.Month(), rev.DateTime.Day(),
			rev.DateTime.Hour(), rev.DateTime.Minute(), rev.Revision)
	}

	return nil
}

func GetRevision(rev int, config *cfg.Config) error {
	revision, err := state.Get(rev, config)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", *revision)

	return nil
}
