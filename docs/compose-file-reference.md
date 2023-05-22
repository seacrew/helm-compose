# Compose File Reference

## storage

More details regarding the available storage providers and all options can be foudn [here](storage-providers.md).

```yaml
storage:
	name: my-compose
	type: local
```

| Option            | Type   | Description | Required | Default |
| ----------------- | ------ | ----------- | -------- | ------- |
| name              | string |             | true     |         |
| type              | string |             | false    | local   |
| numberOfRevisions | int    |             | false    | 10      |

## releases

```yaml
releases:
	my-website:
		chart: bitnami/wordpress
```

| Option           | Type   | Description | Required | Default        |
| ---------------- | ------ | ----------- | -------- | -------------- |
| name             | string |             | false    |                |
| chart            | string |             | true     |                |
| chartVersion     | string |             | false    | latest         |
| forceUpdate      | bool   |             | false    | false          |
| historyMax       | int    |             | false    | 10             |
| createNamespace  | bool   |             | false    | false          |
| cleanUpOnFail    | bool   |             | false    | false          |
| dependencyUpdate | bool   |             | false    | false          |
| skipTlsVerify    | bool   |             | false    | false          |
| skipCrds         | bool   |             | false    | false          |
| postRenderer     | string |             | false    |                |
| postRendererArgs | array  |             | false    |                |
| kubeconfig       | string |             | false    | ~/.kube/config |
| kubecontext      | string |             | false    |                |
| caFile           | string |             | false    |                |
| certFile         | string |             | false    |                |
| keyFile          | string |             | false    |                |
| timeout          | string |             | false    | 5m             |
| values           | map    |             | false    |                |
| valueFiles       | string |             | false    | 5m             |

Uninstall options:

| Option           | Type   | Description | Required | Default    |
| ---------------- | ------ | ----------- | -------- | ---------- |
| deletionStrategy | string |             | false    | background |
| deletionTimeout  | string |             | false    | 5m         |
| deletionNoHooks  | bool   |             | false    | false      |
| keepHistory      | bool   |             | false    | false      |

## repositories

A map of repository names and their respective urls.

```yaml
repositories:
  bitnami: https://charts.bitnami.com/bitnami
```
