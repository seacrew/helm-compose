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
type Compose struct {
	Name        string `yaml:"name,omitempty"`
	State       string `yaml:"state,omitempty"`
	Namespace   string `yaml:"namespace,omitempty"`
	KubeConfig  string `yaml:"kubeconfig,omitempty"`
	KubeContext string `yaml:"kubecontext,omitempty"`
}

type Config struct {
	Version      string             `yaml:"apiVersion,omitempty"`
	Compose      Compose            `yaml:"compose,omitempty"`
	Releases     map[string]Release `yaml:"releases,omitempty"`
	Repositories map[string]string  `yaml:"repositories,omitempty"`
}
