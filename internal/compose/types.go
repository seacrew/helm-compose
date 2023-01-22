package compose

type Release struct {
	Name            string                      `yaml:"name,omitempty"`
	Chart           string                      `yaml:"chart,omitempty"`
	ChartVersion    string                      `yaml:"chartVersion,omitempty"`
	Namespace       string                      `yaml:"namespace,omitempty"`
	CreateNamespace bool                        `yaml:"createNamespace,omitempty"`
	KubeConfig      string                      `yaml:"kubeconfig,omitempty"`
	KubeContext     string                      `yaml:"kubecontext,omitempty"`
	Values          map[interface{}]interface{} `yaml:"values,omitempty"`
	ValueFiles      []string                    `yaml:"valueFiles,omitempty"`
}

type Config struct {
	Version      string             `yaml:"composeVersion,omitempty"`
	Name         string             `yaml:"name,omitempty"`
	State        string             `yaml:"state,omitempty"`
	Releases     map[string]Release `yaml:"releases,omitempty"`
	Repositories map[string]string  `yaml:"repositories,omitempty"`
}
