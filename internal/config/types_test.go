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
package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigEqualVersion(t *testing.T) {
	cfg := Config{Version: "v1.0"}
	otherCfg := Config{Version: "v1.0"}
	assert.True(t, cfg.Equal(&otherCfg))

	cfg = Config{Version: "v1.0"}
	otherCfg = Config{Version: "v1.1"}
	assert.False(t, cfg.Equal(&otherCfg))

}

func TestConfigEqualStorage(t *testing.T) {
	cfg := Config{
		Storage: Storage{
			Type: Local,
			Name: "mycompose",
			Path: ".state",
		},
	}
	assert.True(t, cfg.Equal(&cfg))

	otherCfg := Config{
		Storage: Storage{
			Type: Local,
			Name: "mycompose",
		},
	}
	assert.False(t, cfg.Equal(&otherCfg))
}

func TestConfigEqualRepositories(t *testing.T) {
	cfg := Config{
		Repositories: map[string]string{},
	}

	otherCfg := Config{
		Repositories: map[string]string{
			"bitnami": "https://charts.bitnami.com/bitnami",
		},
	}
	assert.False(t, cfg.Equal(&otherCfg))

	cfg = Config{
		Repositories: map[string]string{
			"bitnami": "https://charts.bitnami.com/bitnami",
		},
	}

	otherCfg = Config{
		Repositories: map[string]string{
			"bitnami": "https://charts.bitnami.com/bitnami",
		},
	}
	assert.True(t, cfg.Equal(&otherCfg))

	cfg = Config{
		Repositories: map[string]string{
			"bitnami": "https://charts.bitnami.com/bitnami",
		},
	}

	otherCfg = Config{
		Repositories: map[string]string{
			"bitnami": "https://charts.bitnami.net/bitnami",
		},
	}
	assert.False(t, cfg.Equal(&otherCfg))
}

func TestConfigEqualReleases(t *testing.T) {
	cfg := Config{
		Releases: map[string]Release{},
	}

	otherCfg := Config{
		Releases: map[string]Release{
			"wordpress": {
				Chart:        "bitnami/wordpress",
				ChartVersion: "14.2.1",
			},
		},
	}
	assert.False(t, cfg.Equal(&otherCfg))

	cfg = Config{
		Releases: map[string]Release{
			"wordpress": {
				Chart:        "bitnami/wordpress",
				ChartVersion: "14.2.1",
				Values: map[string]interface{}{
					"wordpressBlogName": "my-site",
				},
			},
		},
	}

	otherCfg = Config{
		Releases: map[string]Release{
			"wordpress": {
				Chart:        "bitnami/wordpress",
				ChartVersion: "14.2.1",
				Values: map[string]interface{}{
					"wordpressBlogName": "my-site",
				},
			},
		},
	}
	assert.True(t, cfg.Equal(&otherCfg))

	cfg = Config{
		Releases: map[string]Release{
			"wordpress": {
				Chart:        "bitnami/wordpress",
				ChartVersion: "14.2.1",
				Values: map[string]interface{}{
					"wordpressBlogName": "my-site",
				},
			},
		},
	}

	otherCfg = Config{
		Releases: map[string]Release{
			"wordpress": {
				Chart:        "bitnami/wordpress",
				ChartVersion: "14.2.1",
				Values: map[string]interface{}{
					"wordpressBlogName": "my-new-site",
				},
			},
		},
	}
	assert.False(t, cfg.Equal(&otherCfg))
}
