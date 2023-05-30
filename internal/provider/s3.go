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
	"crypto/tls"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	cfg "github.com/seacrew/helm-compose/internal/config"
	"github.com/seacrew/helm-compose/internal/util"
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
	return nil, nil
}

func (p S3Provider) store(encodedConfig *string) error {
	return nil
}
func (p S3Provider) list() ([]ComposeRevision, error) {
	return nil, nil
}
func (p S3Provider) get(revision int) (*[]byte, error) {
	return nil, nil
}
