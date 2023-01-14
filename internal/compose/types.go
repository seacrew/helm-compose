package compose

type Release struct {
	Name         string                 `yaml:"name,omitempty"`
	Chart        string                 `yaml:"chart,omitempty"`
	ChartVersion string                 `yaml:"chartVersion,omitempty"`
	Namespace    string                 `yaml:"namespace,omitempty"`
	KubeContext  string                 `yaml:"kubeContext,omitempty"`
	Values       map[string]interface{} `yaml:"values,omitempty"`
	ValueFiles   []string               `yaml:"valueFiles,omitempty"`
}

type Config struct {
	Version      string             `yaml:"composeVersion,omitempty"`
	Releases     map[string]Release `yaml:"releases,omitempty"`
	Repositories map[string]string  `yaml:"repositories,omitempty"`
}
