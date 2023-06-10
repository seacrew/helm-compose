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
	"sort"
	"time"

	cfg "github.com/seacrew/helm-compose/internal/config"
	"github.com/seacrew/helm-compose/internal/util"
	"gopkg.in/yaml.v2"
)

type ComposeRevision struct {
	Revision int
	DateTime time.Time
}

type Provider interface {
	load() (*[]byte, error)
	store(encodedConfig *string) error
	list() ([]ComposeRevision, error)
	get(revision int) (*[]byte, error)
}

var provider Provider

func getProvider(providerConfig *cfg.Storage) (Provider, error) {
	if provider != nil {
		return provider, nil
	}

	if providerConfig.NumberOfRevisions <= 0 {
		providerConfig.NumberOfRevisions = 10
	}

	var err error

	switch providerConfig.Type {
	case cfg.Local:
		provider = newLocalProvider(providerConfig)
		return provider, nil
	case cfg.Kubernetes:
		provider, err = newKubernetesProvider(providerConfig)
		return provider, err
	case cfg.S3:
		provider, err = newS3Provider(providerConfig)
		return provider, err
	default:
		return nil, fmt.Errorf("unknown provider type %q", providerConfig.Type)
	}
}

func Load(config *cfg.Config) (*cfg.Config, error) {
	provider, err := getProvider(&config.Storage)
	if err != nil {
		return nil, err
	}

	data, err := provider.load()
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	prevConfig, err := util.DecodeComposeConfig(string(*data))
	if err != nil {
		return nil, err
	}

	if err := cfg.ValidateCompose(prevConfig); err != nil {
		return nil, fmt.Errorf("couldn't load previous revision %s", err.Error())
	}

	return prevConfig, nil
}

func Store(config *cfg.Config) error {
	provider, err := getProvider(&config.Storage)
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

func List(config *cfg.Config) ([]ComposeRevision, error) {
	provider, err := getProvider(&config.Storage)
	if err != nil {
		return nil, err
	}

	revisions, err := provider.list()
	if err != nil {
		return nil, err
	}

	sort.Slice(revisions, func(i, j int) bool {
		return revisions[i].Revision < revisions[j].Revision
	})

	return revisions, nil
}

func Get(revision int, config *cfg.Config) (*string, error) {
	provider, err := getProvider(&config.Storage)
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
