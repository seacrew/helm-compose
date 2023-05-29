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
	"os"
	"regexp"
	"strconv"

	cfg "github.com/seacrew/helm-compose/internal/config"
	"github.com/seacrew/helm-compose/internal/util"
)

const (
	pathFormat = "%s/%s-%d"
)

type LocalProvider struct {
	name              string
	path              string
	numberOfRevisions int
}

func newLocalProvider(providerConfig *cfg.Storage) *LocalProvider {
	provider := LocalProvider{
		name:              providerConfig.Name,
		path:              providerConfig.Path,
		numberOfRevisions: providerConfig.NumberOfRevisions,
	}

	if len(provider.path) == 0 {
		provider.path = ".hcstate"
	}

	return &provider
}

func (p LocalProvider) load() (*[]byte, error) {
	if _, err := os.Stat(p.path); os.IsNotExist(err) {
		if err := os.Mkdir(p.path, os.ModePerm); err != nil {
			return nil, err
		}
	}

	_, maximum, err := p.minMax(p.name, p.path)
	if err != nil {
		return nil, err
	}

	if maximum == 0 {
		return nil, nil
	}

	file, err := os.ReadFile(fmt.Sprintf(pathFormat, p.path, p.name, maximum))
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (p LocalProvider) store(encodedConfig *string) error {
	minimum, maximum, err := p.minMax(p.name, p.path)
	if err != nil {
		return err
	}

	maximum = maximum + 1

	if err := os.WriteFile(fmt.Sprintf(pathFormat, p.path, p.name, maximum), []byte(*encodedConfig), 0644); err != nil {
		return err
	}

	if minimum > maximum-p.numberOfRevisions {
		return nil
	}

	for i := minimum; i <= maximum-p.numberOfRevisions; i++ {
		if err := os.Remove(fmt.Sprintf(pathFormat, p.path, p.name, i)); err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func (p LocalProvider) list() ([]ComposeRevision, error) {
	files, err := os.ReadDir(p.path)
	if err != nil {
		return nil, err
	}

	r, _ := regexp.Compile(fmt.Sprintf("^%s-(\\d+)$", p.name))

	revisions := []ComposeRevision{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		matches := r.FindStringSubmatch(file.Name())
		if len(matches) == 0 {
			continue
		}

		revision, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, err
		}

		info, err := file.Info()
		if err != nil {
			return nil, err
		}

		revisions = append(revisions, ComposeRevision{revision, info.ModTime()})
	}

	return revisions, nil
}

func (p LocalProvider) get(revision int) (*[]byte, error) {
	file, err := os.ReadFile(fmt.Sprintf(pathFormat, p.path, p.name, revision))
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (p LocalProvider) minMax(name string, path string) (int, int, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return -1, -1, err
	}

	r, _ := regexp.Compile(fmt.Sprintf("^%s-(\\d+)$", name))

	revisions := []int{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		matches := r.FindStringSubmatch(file.Name())
		if len(matches) == 0 {
			continue
		}

		revision, err := strconv.Atoi(matches[1])
		if err != nil {
			return -1, -1, nil
		}

		revisions = append(revisions, revision)
	}

	minimum, maximum := util.MinMax(revisions)
	return minimum, maximum, nil
}
