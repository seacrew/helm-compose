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

type Config struct {
	Version      string             `yaml:"apiVersion,omitempty"`
	State        State              `yaml:"state,omitempty"`
	Releases     map[string]Release `yaml:"releases,omitempty"`
	Repositories map[string]string  `yaml:"repositories,omitempty"`
}

type Release struct {
	Name             string                 `yaml:"name,omitempty"`
	Chart            string                 `yaml:"chart,omitempty"`
	ChartVersion     string                 `yaml:"chartVersion,omitempty"`
	Namespace        string                 `yaml:"namespace,omitempty"`
	ForceUpdate      bool                   `yaml:"forceUpdate,omitempty"`
	HistoryMax       int                    `yaml:"historyMax,omitempty"`
	CreateNamespace  bool                   `yaml:"createNamespace,omitempty"`
	CleanUpOnFail    bool                   `yaml:"cleanupOnFail,omitempty"`
	DependencyUpdate bool                   `yaml:"dependencyUpdate,omitempty"`
	SkipTLSVerify    bool                   `yaml:"skipTlsVerify,omitempty"`
	SkipCRDs         bool                   `yaml:"skipCrds,omitempty"`
	PostRenderer     string                 `yaml:"postRenderer,omitempty"`
	PostRendererArgs []string               `yaml:"postRendererArgs,omitempty"`
	KubeConfig       string                 `yaml:"kubeconfig,omitempty"`
	KubeContext      string                 `yaml:"kubecontext,omitempty"`
	CAFile           string                 `yaml:"caFile,omitempty"`
	CertFile         string                 `yaml:"certFile,omitempty"`
	KeyFile          string                 `yaml:"keyFile,omitempty"`
	Timeout          string                 `yaml:"timeout,omitempty"`
	Values           map[string]interface{} `yaml:"values,omitempty"`
	ValueFiles       []string               `yaml:"valueFiles,omitempty"`

	// Uninstall flags
	DeletionStrategy string `yaml:"deletionStrategy,omitempty"`
	DeletionTimeout  string `yaml:"deletionTimeout,omitempty"`
	DeletionNoHooks  bool   `yaml:"deletionNoHooks,omitempty"`
	KeepHistory      bool   `yaml:"keepHistory,omitempty"`
}

type State struct {
	Type           ProviderType `yaml:"type,omitempty"`
	Name           string       `yaml:"name,omitempty"`
	NumberOfStates int          `yaml:"numberOfStates,omitempty"`
	// Local storage fields
	Path string `yaml:"path,omitempty"`
	// K8s storage fields
	Namespace string `yaml:"namespace,omitempty"`
}

type ProviderType string

const (
	Local      ProviderType = "local"
	Kubernetes ProviderType = "kubernetes"
	S3         ProviderType = "s3"
)
