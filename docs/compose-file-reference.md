# Compose File Reference

## storage

```yaml
storage:
  name: my-compose
  type: local
  numberOfRevisions: 10
```

| Option            | Type   | Description                                                                                 | Required | Default | ApiVersion |
| ----------------- | ------ | ------------------------------------------------------------------------------------------- | -------- | ------- | ---------- |
| name              | string | Name to be used to store revisions with a storage provider.                                 | true     |         | 1.0        |
| type              | string | Type / name of the storage provider you want to use. By default local files will be stored. | false    | local   | 1.0        |
| numberOfRevisions | int    | Number of revisions to be stored and to be able to rollback to.                             | false    | 10      | 1.0        |

More details regarding the available storage providers and provider specific options can be found [here.](storage-providers.md)

## releases

```yaml
releases:
  my-website:
    chart: bitnami/wordpress
```

| Option           | Type   | Description                                                                                                                                                                 | Required | Default        | ApiVersion |
| ---------------- | ------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | -------------- | ---------- |
| chart            | string | Name of the chart to be used.                                                                                                                                               | true     |                | 1.0        |
| chartVersion     | string | Version of the chart to be used.                                                                                                                                            | false    | latest         | 1.0        |
| forceUpdate      | bool   | Force resource updates through a replacement strategy                                                                                                                       | false    | false          | 1.0        |
| historyMax       | int    | Limit the maximum number of revisions saved per release.                                                                                                                    | false    | 10             | 1.0        |
| createNamespace  | bool   | Create the release namespace if not present                                                                                                                                 | false    | false          | 1.0        |
| cleanUpOnFail    | bool   | Allow deletion of new resources created in this upgrade when upgrade fails.                                                                                                 | false    | false          | 1.0        |
| dependencyUpdate | bool   | Update dependencies if they are missing before installing the chart                                                                                                         | false    | false          | 1.0        |
| skipTlsVerify    | bool   | Skip tls certificate checks for the chart download                                                                                                                          | false    | false          | 1.0        |
| skipCrds         | bool   | If set, no CRDs will be installed.                                                                                                                                          | false    | false          | 1.0        |
| postRenderer     | string | The path to an executable to be used for post rendering. If it exists in $PATH, the binary will be used, otherwise it will try to look for the executable at the given path | false    |                | 1.0        |
| postRendererArgs | array  | An argument to the post-renderer (can specify multiple) (default [])                                                                                                        | false    |                | 1.0        |
| kubeconfig       | string | Path to the kubeconfig file                                                                                                                                                 | false    | ~/.kube/config | 1.0        |
| kubecontext      | string | Name of the kubeconfig context to use                                                                                                                                       | false    |                | 1.0        |
| caFile           | string | Verify certificates of HTTPS-enabled servers using this CA bundle                                                                                                           | false    |                | 1.0        |
| certFile         | string | Identify HTTPS client using this SSL certificate file                                                                                                                       | false    |                | 1.0        |
| keyFile          | string | Identify HTTPS client using this SSL key file                                                                                                                               | false    |                | 1.0        |
| timeout          | string | Time to wait for any individual Kubernetes operation (like Jobs for hooks) (default 5m0s)                                                                                   | false    | 5m             | 1.0        |
| wait             | bool   | Waits until all Pods are in a ready state,  It will wait for as long as the --timeout value                                                                                 | false    |                | 1.1        |
| values           | map    | Map of values with highest priority to overwrite any values in the chart values or your additional values files. (Allows for usage of environment variables.)               | false    |                | 1.0        |
| valueFiles       | string | List of paths to value files.                                                                                                                                               | false    | 5m             | 1.0        |

Uninstall options:

| Option           | Type   | Description                                                                                                  | Required | Default    | ApiVersion |
| ---------------- | ------ | ------------------------------------------------------------------------------------------------------------ | -------- | ---------- | ---------- |
| deletionStrategy | string | Must be "background", "orphan", or "foreground". Selects the deletion cascading strategy for the dependents. | false    | background | 1.0        |
| deletionTimeout  | string | Time to wait for any individual Kubernetes operation (like Jobs for hooks) (default 5m0s)                    | false    | 5m         | 1.0        |
| deletionNoHooks  | bool   | Prevent hooks from running during uninstallation                                                             | false    | false      | 1.0        |
| keepHistory      | bool   | Remove all associated resources and mark the release as deleted, but retain the release history              | false    | false      | 1.0        |

## repositories

A map of repository names and their respective urls.

```yaml
repositories:
  bitnami: https://charts.bitnami.com/bitnami
```
