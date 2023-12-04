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
	"bytes"
	"context"
	"fmt"
	"math"
	"regexp"
	"strconv"

	"cloud.google.com/go/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	cfg "github.com/seacrew/helm-compose/internal/config"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const (
	gcsObjectNameFormat  = "%s.v%d.hcstate"
	gcsObjectNamePattern = "%s.v(\\d+).hcstate$"
)

type GCSProvider struct {
	name              string
	numberOfRevisions int
	bucket            *string
	prefix            *string
	client            *storage.Client
}

func newGCSProvider(providerConfig *cfg.Storage) (*GCSProvider, error) {
	var client *storage.Client
	var err error

	if providerConfig.GCSCredentialsFile != "" {
		client, err = storage.NewClient(context.Background(), option.WithCredentialsFile(providerConfig.GCSCredentialsFile))
	} else {
		client, err = storage.NewClient(context.Background())
	}

	if err != nil {
		return nil, err
	}

	return &GCSProvider{
		name:              providerConfig.Name,
		numberOfRevisions: providerConfig.NumberOfRevisions,
		bucket:            &providerConfig.GCSBucket,
		prefix:            &providerConfig.GCSPrefix,
		client:            client,
	}, nil
}

func (p GCSProvider) load() (*[]byte, error) {
	it := p.client.Bucket(*p.bucket).Objects(context.Background(), &storage.Query{Prefix: *p.prefix})

	_, _, latest, err := p.minMax(it)
	if err != nil {
		return nil, err
	}

	rc, err := p.client.Bucket(*p.bucket).Object(latest.Name).NewReader(context.Background())
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	buff := &bytes.Buffer{}
	if _, err := buff.ReadFrom(rc); err != nil {
		return nil, err
	}

	content := buff.Bytes()
	return &content, nil
}

func (p GCSProvider) store(encodedConfig *string) error {
	it := p.client.Bucket(*p.bucket).Objects(context.Background(), &storage.Query{Prefix: *p.prefix})

	minimum, maximum, _, err := p.minMax(it)
	if err != nil {
		return err
	}

	revision := maximum + 1

	key := fmt.Sprintf(gcsObjectNameFormat, p.name, revision)
	if len(*p.prefix) > 0 {
		key = fmt.Sprintf("%s/%s", *p.prefix, key)
	}

	buf := []byte(*encodedConfig)

	wc := p.client.Bucket(*p.bucket).Object(key).NewWriter(context.Background())
	if _, err = wc.Write(buf); err != nil {
		return err
	}

	if minimum > revision-p.numberOfRevisions {
		return nil
	}

	for i := minimum; i < revision-p.numberOfRevisions; i++ {
		key := fmt.Sprintf(gcsObjectNameFormat, p.name, i)
		if len(*p.prefix) > 0 {
			key = fmt.Sprintf("%s/%s", *p.prefix, key)
		}

		o := p.client.Bucket(*p.bucket).Object(key)
		if err := o.Delete(context.Background()); err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func (p GCSProvider) list() ([]ComposeRevision, error) {
	resp, err := p.lister.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: p.bucket, Prefix: p.prefix})
	if err != nil {
		return nil, err
	}

	var revisions []ComposeRevision

	r, err := regexp.Compile(fmt.Sprintf(s3ObjectNamePattern, p.name))
	if err != nil {
		return nil, err
	}

	for _, item := range resp.Contents {
		matches := r.FindStringSubmatch(*item.Key)
		if len(matches) == 0 {
			continue
		}

		revision, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, err
		}

		revisions = append(revisions, ComposeRevision{revision, *item.LastModified})
	}

	return revisions, nil
}

func (p GCSProvider) get(revision int) (*[]byte, error) {
	key := fmt.Sprintf(s3ObjectNameFormat, p.name, revision)
	if len(*p.prefix) > 0 {
		key = fmt.Sprintf("%s/%s", *p.prefix, key)
	}

	buff := &aws.WriteAtBuffer{}
	_, err := p.downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: p.bucket,
			Key:    &key,
		})

	if err != nil {
		return nil, err
	}

	content := buff.Bytes()
	return &content, nil
}

func (p GCSProvider) minMax(it *storage.ObjectIterator) (int, int, *storage.ObjectAttrs, error) {
	minimum, maximum := math.MaxInt, 0
	var latest *storage.ObjectAttrs

	r, err := regexp.Compile(fmt.Sprintf(gcsObjectNamePattern, p.name))
	if err != nil {
		return -1, -1, nil, err
	}

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}

		matches := r.FindStringSubmatch(attrs.Name)
		if len(matches) == 0 {
			continue
		}

		revision, err := strconv.Atoi(matches[1])
		if err != nil {
			return -1, -1, nil, err
		}

		if revision < minimum {
			minimum = revision
		}

		if revision > maximum {
			maximum = revision
			latest = attrs
		}
	}

	return minimum, maximum, latest, nil
}
