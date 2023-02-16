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
	ValueFiles      []string               `yaml:"valuefiles,omitempty"`
}

type Config struct {
	Version string `yaml:"apiVersion,omitempty"`
	State   struct {
		Name    string `yaml:"name,omitempty"`
		Storage string `yaml:"storage,omitempty"`
	} `yaml:"state,omitempty"`
	Releases     map[string]Release `yaml:"releases,omitempty"`
	Repositories map[string]string  `yaml:"repositories,omitempty"`
}
