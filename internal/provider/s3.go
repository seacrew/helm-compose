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
	"crypto/tls"
	"fmt"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	cfg "github.com/seacrew/helm-compose/internal/config"
	"github.com/seacrew/helm-compose/internal/util"
)

const (
	s3ObjectNameFormat  = "%s.v%d.hcstate"
	s3ObjectNamePattern = "%s.v(\\d+).hcstate$"
)

type S3Provider struct {
	name              string
	numberOfRevisions int
	bucket            *string
	prefix            *string
	lister            *s3.S3
	uploader          *s3manager.Uploader
	downloader        *s3manager.Downloader
}

func newS3Provider(providerConfig *cfg.Storage) (*S3Provider, error) {
	config := &aws.Config{}

	if len(providerConfig.S3Region) > 0 {
		config.Region = &providerConfig.S3Region
	} else if os.Getenv("AWS_REGION") != "" {
		config.Region = aws.String(os.Getenv("AWS_REGION"))
	} else {
		return nil, fmt.Errorf("AWS region not specified")
	}

	if len(providerConfig.S3Endpoint) > 0 {
		config.Endpoint = &providerConfig.S3Endpoint
	}

	if providerConfig.S3DisableSSL {
		config.DisableSSL = util.NewBool(true)
	}

	if providerConfig.S3Insecure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		config.HTTPClient = &http.Client{Transport: tr}
	}

	if providerConfig.S3ForcePathStyle {
		config.S3ForcePathStyle = util.NewBool(true)
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	provider := &S3Provider{
		name:              providerConfig.Name,
		numberOfRevisions: providerConfig.NumberOfRevisions,
		bucket:            &providerConfig.S3Bucket,
		prefix:            &providerConfig.S3Prefix,
		lister:            s3.New(sess),
		uploader:          s3manager.NewUploader(sess),
		downloader:        s3manager.NewDownloader(sess),
	}

	return provider, nil
}

func (p S3Provider) load() (*[]byte, error) {
	resp, err := p.lister.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: p.bucket, Prefix: p.prefix})
	if err != nil {
		return nil, err
	}

	if len(resp.Contents) == 0 {
		return nil, nil
	}

	_, _, latest, err := p.minMax(resp.Contents)
	if err != nil {
		return nil, err
	}

	if latest == nil {
		return nil, nil
	}

	buff := &aws.WriteAtBuffer{}
	if _, err := p.downloader.Download(buff, &s3.GetObjectInput{Bucket: p.bucket, Key: latest.Key}); err != nil {
		return nil, err
	}

	content := buff.Bytes()
	return &content, nil
}

func (p S3Provider) store(encodedConfig *string) error {
	resp, err := p.lister.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: p.bucket, Prefix: p.prefix})
	if err != nil {
		return err
	}

	minimum, maximum, _, err := p.minMax(resp.Contents)
	if err != nil {
		return err
	}

	revision := maximum + 1

	key := fmt.Sprintf(s3ObjectNameFormat, p.name, revision)
	if len(*p.prefix) > 0 {
		key = fmt.Sprintf("%s/%s", *p.prefix, key)
	}

	reader := bytes.NewReader([]byte(*encodedConfig))

	_, err = p.uploader.Upload(
		&s3manager.UploadInput{
			Bucket: p.bucket,
			Key:    &key,
			Body:   reader,
		})

	if err != nil {
		return err
	}

	if minimum > revision-p.numberOfRevisions {
		return nil
	}

	for i := minimum; i < revision-p.numberOfRevisions; i++ {
		key := fmt.Sprintf(s3ObjectNameFormat, p.name, i)
		if len(*p.prefix) > 0 {
			key = fmt.Sprintf("%s/%s", *p.prefix, key)
		}

		_, err := p.lister.DeleteObject(&s3.DeleteObjectInput{Bucket: p.bucket, Key: &key})
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func (p S3Provider) list() ([]ComposeRevision, error) {
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

func (p S3Provider) get(revision int) (*[]byte, error) {
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

func (p S3Provider) minMax(objects []*s3.Object) (int, int, *s3.Object, error) {
	if len(objects) == 0 {
		return 0, 0, nil, nil
	}

	minimum, maximum := math.MaxInt, 0
	var latest *s3.Object

	r, err := regexp.Compile(fmt.Sprintf(s3ObjectNamePattern, p.name))
	if err != nil {
		return -1, -1, nil, err
	}

	for _, item := range objects {
		matches := r.FindStringSubmatch(*item.Key)

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
			latest = item
		}
	}

	return minimum, maximum, latest, nil
}
