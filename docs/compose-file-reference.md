# Compose File Reference

```go
type Config struct {
	Version      string             `yaml:"apiVersion,omitempty"`
	Storage      Storage            `yaml:"storage,omitempty"`
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

type Storage struct {
	Type              ProviderType `yaml:"type,omitempty"`
	Name              string       `yaml:"name,omitempty"`
	NumberOfRevisions int          `yaml:"numberOfRevisions,omitempty"`
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
```
