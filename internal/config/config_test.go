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
	"log"
	"testing"
)

func TestParseSimpleConfig(t *testing.T) {
	config, err := ParseComposeFile("../../examples/simple-compose.yaml")

	if err != nil {
		log.Fatal(err)
	}

	if config.Storage.Name != "simple" {
		log.Fatalf("Was expecting revision name 'simple' but got '%s'", config.Storage.Name)
	}

	if config.Storage.Type != Local {
		log.Fatalf("Was expecting revision provider type '%s' but got '%s'", Local, config.Storage.Type)
	}

	if len(config.Releases) != 2 {
		log.Fatalf("Was expecting 2 release but got %d", len(config.Releases))
	}
}
