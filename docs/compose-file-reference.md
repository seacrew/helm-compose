# Compose File Reference

## storage

```yaml
storage:
  name: my-compose
  type: local
  numberOfRevisions: 10
```

| Option            | Type   | Description                                                                                 | Required | Default |
| ----------------- | ------ | ------------------------------------------------------------------------------------------- | -------- | ------- |
| name              | string | Name to be used to store revisions with a storage provider.                                 | true     |         |
| type              | string | Type / name of the storage provider you want to use. By default local files will be stored. | false    | local   |
| numberOfRevisions | int    | Number of revisions to be stored and to be able to rollback to.                             | false    | 10      |

More details regarding the available storage providers and provider specific options can be found [here.](storage-providers.md)

## releases

```yaml
releases:
  my-website:
    chart: bitnami/wordpress
```

| Option           | Type   | Description                                                                                                                                                                 | Required | Default        |
| ---------------- | ------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | -------------- |
| chart            | string | Name of the chart to be used.                                                                                                                                               | true     |                |
| chartVersion     | string | Version of the chart to be used.                                                                                                                                            | false    | latest         |
| forceUpdate      | bool   | Force resource updates through a replacement strategy                                                                                                                       | false    | false          |
| historyMax       | int    | Limit the maximum number of revisions saved per release.                                                                                                                    | false    | 10             |
| createNamespace  | bool   | Create the release namespace if not present                                                                                                                                 | false    | false          |
| cleanUpOnFail    | bool   | Allow deletion of new resources created in this upgrade when upgrade fails.                                                                                                 | false    | false          |
| dependencyUpdate | bool   | Update dependencies if they are missing before installing the chart                                                                                                         | false    | false          |
| skipTlsVerify    | bool   | Skip tls certificate checks for the chart download                                                                                                                          | false    | false          |
| skipCrds         | bool   | If set, no CRDs will be installed.                                                                                                                                          | false    | false          |
| postRenderer     | string | The path to an executable to be used for post rendering. If it exists in $PATH, the binary will be used, otherwise it will try to look for the executable at the given path | false    |                |
| postRendererArgs | array  | An argument to the post-renderer (can specify multiple) (default [])                                                                                                        | false    |                |
| kubeconfig       | string | Path to the kubeconfig file                                                                                                                                                 | false    | ~/.kube/config |
| kubecontext      | string | Name of the kubeconfig context to use                                                                                                                                       | false    |                |
| caFile           | string | Verify certificates of HTTPS-enabled servers using this CA bundle                                                                                                           | false    |                |
| certFile         | string | Identify HTTPS client using this SSL certificate file                                                                                                                       | false    |                |
| keyFile          | string | Identify HTTPS client using this SSL key file                                                                                                                               | false    |                |
| timeout          | string | Time to wait for any individual Kubernetes operation (like Jobs for hooks) (default 5m0s)                                                                                   | false    | 5m             |
| wait             | bool   | Waits until all Pods are in a ready state,  It will wait for as long as the --timeout value                                                                                 | false    |

| values           | map    | Map of values with highest priority to overwrite any values in the chart values or your additional values files. (Allows for usage of environment variables.)               | false    |                |
| valueFiles       | string | List of paths to value files.                                                                                                                                               | false    | 5m             |

Uninstall options:

| Option           | Type   | Description                                                                                                  | Required | Default    |
| ---------------- | ------ | ------------------------------------------------------------------------------------------------------------ | -------- | ---------- |
| deletionStrategy | string | Must be "background", "orphan", or "foreground". Selects the deletion cascading strategy for the dependents. | false    | background |
| deletionTimeout  | string | Time to wait for any individual Kubernetes operation (like Jobs for hooks) (default 5m0s)                    | false    | 5m         |
| deletionNoHooks  | bool   | Prevent hooks from running during uninstallation                                                             | false    | false      |
| keepHistory      | bool   | Remove all associated resources and mark the release as deleted, but retain the release history              | false    | false      |

## repositories

A map of repository names and their respective urls.

```yaml
repositories:
  bitnami: https://charts.bitnami.com/bitnami
```
