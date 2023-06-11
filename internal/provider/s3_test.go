package provider

import (
	"log"
	"testing"

	cfg "github.com/seacrew/helm-compose/internal/config"
)

func TestS3List(t *testing.T) {
	config := &cfg.Storage{
		Name:              "test",
		NumberOfRevisions: 5,
		S3Bucket:          "helm-compose",
		S3Prefix:          "test",
		S3Region:          "eu-central-1",
	}

	provider, err := newS3Provider(config)

	if err != nil {
		log.Fatal(err)
	}

	revisions, err := provider.list()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(revisions)
}
