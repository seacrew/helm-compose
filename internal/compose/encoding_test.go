package compose

import (
	"log"
	"testing"
)

func TestEncoding(t *testing.T) {
	config := Config{
		Version: "1.0",
		Compose: Compose{
			Name:  "mycompose",
			State: "local",
		},
		Releases: map[string]Release{
			"myrelease": {
				Chart: "mychart",
			},
		},
	}

	enc, err := encodeComposeConfig(&config)

	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(decodeComposeConfig(enc))
}
