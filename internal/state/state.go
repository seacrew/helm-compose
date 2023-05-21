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
package state

import (
	cfg "github.com/seacrew/helm-compose/internal/config"
	prov "github.com/seacrew/helm-compose/internal/provider"
	"gopkg.in/yaml.v2"
)

func Load(config *cfg.Config) (*cfg.Config, error) {
	provider, err := NewProvider(&config.State)
	if err != nil {
		return nil, err
	}

	data, err := provider.Load()

	if data == nil {
		return nil, nil
	}

	prevConfig, err := decodeComposeConfig(string(*data))
	if err != nil {
		return nil, err
	}

	return prevConfig, nil
}

func Store(config *cfg.Config) error {
	provider, err := NewProvider(&config.State)
	if err != nil {
		return err
	}

	encodedConfig, err := encodeComposeConfig(config)

	if err != nil {
		return err
	}

	if err := provider.Store(&encodedConfig); err != nil {
		return err
	}

	return nil
}

func List(config *cfg.Config) ([]prov.ReleaseRevision, error) {
	provider, err := NewProvider(&config.State)
	if err != nil {
		return nil, err
	}

	revisions, err := provider.List()
	if err != nil {
		return nil, err
	}

	return revisions, nil
}

func Get(revision int, config *cfg.Config) (*string, error) {
	provider, err := NewProvider(&config.State)
	if err != nil {
		return nil, err
	}

	data, err := provider.Get(revision)
	if err != nil {
		return nil, err
	}

	revConfig, err := decodeComposeConfig(string(*data))
	if err != nil {
		return nil, err
	}

	b, err := yaml.Marshal(revConfig)
	if err != nil {
		return nil, err
	}

	revYaml := string(b)
	return &revYaml, nil
}
