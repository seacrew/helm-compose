package compose

import (
	"log"
	"testing"
)

func TestParseConfig(t *testing.T) {

	config, err := ParseConfig("../../testdata/helm-compose.yaml")

	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(config)
}
