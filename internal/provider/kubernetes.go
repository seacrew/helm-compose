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
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	b64 "encoding/base64"

	cfg "github.com/seacrew/helm-compose/internal/config"
	"github.com/seacrew/helm-compose/internal/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	secretNameFormat  = "helm.compose.%s.v%d"
	secretNamePattern = "^helm.compose.%s.v(\\d+)$"
)

type KubernetesProvider struct {
	name              string
	numberOfRevisions int
	namespace         string
	client            *kubernetes.Clientset
	listOptions       *metav1.ListOptions
}

func newKubernetesProvider(providerConfig *cfg.Storage) (*KubernetesProvider, error) {
	namespace := providerConfig.Namespace
	if len(namespace) == 0 {
		namespace = "default"
	}

	kubeconfig := providerConfig.KubeConfig
	if len(kubeconfig) == 0 {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		kubeconfig = filepath.Join(
			homedir, ".kube", "config",
		)
	}

	var err error
	var config *rest.Config

	if len(providerConfig.KubeContext) == 0 {
		config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}, &clientcmd.ConfigOverrides{CurrentContext: providerConfig.KubeContext}).ClientConfig()
		if err != nil {
			return nil, err
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}

	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"app.kubernetes.io/managed-by": "Helm-Compose", "helm-compose/name": providerConfig.Name}}
	listOptions := metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
	}

	provider := &KubernetesProvider{
		name:              providerConfig.Name,
		numberOfRevisions: providerConfig.NumberOfRevisions,
		namespace:         namespace,
		client:            clientset,
		listOptions:       &listOptions,
	}

	return provider, nil
}

func (p KubernetesProvider) load() (*[]byte, error) {
	secrets, err := p.client.CoreV1().Secrets(p.namespace).List(context.Background(), *p.listOptions)
	if err != nil {
		return nil, err
	}

	if len(secrets.Items) == 0 {
		return nil, nil
	}

	_, _, latest, err := p.minMax(secrets.Items)
	if err != nil {
		return nil, err
	}

	if latest == nil {
		return nil, nil
	}

	data, err := b64.StdEncoding.DecodeString(string(latest.Data["compose"]))
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (p KubernetesProvider) store(encodedConfig *string) error {
	secrets, err := p.client.CoreV1().Secrets(p.namespace).List(context.Background(), *p.listOptions)
	if err != nil {
		return err
	}

	minimum, maximum, _, err := p.minMax(secrets.Items)
	if err != nil {
		return err
	}

	revision := maximum + 1

	data := b64.StdEncoding.EncodeToString([]byte(*encodedConfig))

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(secretNameFormat, p.name, revision),
			Namespace: p.namespace,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "Helm-Compose",
				"helm-compose/name":            p.name,
			},
		},
		Immutable: util.NewBool(true),
		Data: map[string][]byte{
			"compose": []byte(data),
		},
		Type: "helm-compose/revision.v1",
	}

	_, err = p.client.CoreV1().Secrets(p.namespace).Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	if minimum > revision-p.numberOfRevisions {
		return nil
	}

	for i := minimum; i <= revision-p.numberOfRevisions; i++ {
		if err = p.client.CoreV1().Secrets(p.namespace).Delete(context.Background(), fmt.Sprintf(secretNameFormat, p.name, i), metav1.DeleteOptions{}); err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func (p KubernetesProvider) list() ([]ComposeRevision, error) {
	secrets, err := p.client.CoreV1().Secrets(p.namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	if len(secrets.Items) == 0 {
		return nil, nil
	}

	revisions := []ComposeRevision{}

	r, err := regexp.Compile(fmt.Sprintf(secretNamePattern, p.name))
	if err != nil {
		return nil, err
	}

	for _, item := range secrets.Items {
		matches := r.FindStringSubmatch(item.Name)
		if len(matches) == 0 {
			continue
		}

		revision, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, err
		}

		revisions = append(revisions, ComposeRevision{revision, item.CreationTimestamp.Time})
	}

	return revisions, nil
}

func (p KubernetesProvider) get(revision int) (*[]byte, error) {
	secret, err := p.client.CoreV1().Secrets(p.namespace).Get(context.Background(), fmt.Sprintf(secretNameFormat, p.name, revision), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	data, err := b64.StdEncoding.DecodeString(string(secret.Data["compose"]))
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (p KubernetesProvider) minMax(secrets []corev1.Secret) (int, int, *corev1.Secret, error) {
	if len(secrets) == 0 {
		return 0, 0, nil, nil
	}

	minimum, maximum := math.MaxInt, 0

	r, err := regexp.Compile(fmt.Sprintf(secretNamePattern, p.name))
	if err != nil {
		return -1, -1, nil, err
	}

	var latest corev1.Secret

	for _, secret := range secrets {
		matches := r.FindStringSubmatch(secret.Name)
		if len(matches) == 0 {
			continue
		}

		revision, err := strconv.Atoi(matches[1])
		if err != nil {
			return -1, -1, nil, err
		}

		if revision > maximum {
			maximum = revision
			latest = secret
		}

		if revision < minimum {
			minimum = revision
		}
	}

	return minimum, maximum, &latest, nil
}
