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
package compose

type Release struct {
	Name            string                 `yaml:"name,omitempty"`
	Chart           string                 `yaml:"chart,omitempty"`
	ChartVersion    string                 `yaml:"chartVersion,omitempty"`
	Namespace       string                 `yaml:"namespace,omitempty"`
	CreateNamespace bool                   `yaml:"createNamespace,omitempty"`
	KubeConfig      string                 `yaml:"kubeconfig,omitempty"`
	KubeContext     string                 `yaml:"kubecontext,omitempty"`
	Values          map[string]interface{} `yaml:"values,omitempty"`
	ValueFiles      []string               `yaml:"valueFiles,omitempty"`
}

type StorageLocal struct {
}

type StorageKubernetes struct {
}

type Storage struct {
	Type           string `yaml:"type,omitempty"`
	NameOverride   string `yaml:"nameOverride,omitempty"`
	NumberOfStates int    `yaml:"numberOfStates,omitempty"`
	// Local storage fields
	Path string `yaml:"path,omitempty"`
	// K8s storage fields
	Namespace string `yaml:"namespace,omitempty"`
}

type Config struct {
	Version      string             `yaml:"apiVersion,omitempty"`
	Name         string             `yaml:"name,omitempty"`
	Storage      Storage            `yaml:"storage,omitempty"`
	Releases     map[string]Release `yaml:"releases,omitempty"`
	Repositories map[string]string  `yaml:"repositories,omitempty"`
}
