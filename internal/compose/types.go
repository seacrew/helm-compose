package compose

type Release struct {
	Name         string                 `json:"name,omitempty"`
	Chart        string                 `json:"chart,omitempty"`
	ChartVersion string                 `json:"chartVersion,omitempty"`
	Namespace    string                 `json:"namespace,omitempty"`
	KubeContext  string                 `json:"kubeContext,omitempty"`
	Values       map[string]interface{} `json:"values,omitempty"`
	ValueFiles   []string               `json:"valueFiles,omitempty"`
}

type Config struct {
	Version      string             `json:"composeVersion,omitempty"`
	Releases     map[string]Release `json:"releases,omitempty"`
	Repositories map[string]string  `json:"repositories,omitempty"`
}
