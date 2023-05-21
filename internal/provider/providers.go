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
package provider

import (
	"fmt"
	"time"

	cfg "github.com/seacrew/helm-compose/internal/config"
	"github.com/seacrew/helm-compose/internal/util"
	"gopkg.in/yaml.v2"
)

type ReleaseRevision struct {
	Revision int
	DateTime time.Time
}

type Provider interface {
	load() (*[]byte, error)
	store(encodedConfig *string) error
	list() ([]ReleaseRevision, error)
	get(revision int) (*[]byte, error)
}

func newProvider(providerConfig *cfg.Storage) (Provider, error) {
	if providerConfig.NumberOfRevisions <= 0 {
		providerConfig.NumberOfRevisions = 10
	}

	switch providerConfig.Type {
	case cfg.Local:
		return newLocal(providerConfig), nil
	case cfg.Kubernetes:
		return nil, fmt.Errorf("provider type kubernetes is not yet implemented")
	case cfg.S3:
		return nil, fmt.Errorf("provider type s3 is not yet implemented")
	default:
		return nil, fmt.Errorf("unknown provider type %q", providerConfig.Type)
	}
}

func Load(config *cfg.Config) (*cfg.Config, error) {
	provider, err := newProvider(&config.Storage)
	if err != nil {
		return nil, err
	}

	data, err := provider.load()

	if data == nil {
		return nil, nil
	}

	prevConfig, err := util.DecodeComposeConfig(string(*data))
	if err != nil {
		return nil, err
	}

	return prevConfig, nil
}

func Store(config *cfg.Config) error {
	provider, err := newProvider(&config.Storage)
	if err != nil {
		return err
	}

	encodedConfig, err := util.EncodeComposeConfig(config)

	if err != nil {
		return err
	}

	if err := provider.store(&encodedConfig); err != nil {
		return err
	}

	return nil
}

func List(config *cfg.Config) ([]ReleaseRevision, error) {
	provider, err := newProvider(&config.Storage)
	if err != nil {
		return nil, err
	}

	revisions, err := provider.list()
	if err != nil {
		return nil, err
	}

	return revisions, nil
}

func Get(revision int, config *cfg.Config) (*string, error) {
	provider, err := newProvider(&config.Storage)
	if err != nil {
		return nil, err
	}

	data, err := provider.get(revision)
	if err != nil {
		return nil, err
	}

	revConfig, err := util.DecodeComposeConfig(string(*data))
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
