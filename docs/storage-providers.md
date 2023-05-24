# Storage Providers

Following options are applied regardless of the selected provider.

| Option            | Type   | Description | Required | Default |
| ----------------- | ------ | ----------- | -------- | ------- |
| name              | string |             | true     |         |
| type              | string |             | false    | local   |
| numberOfRevisions | int    |             | false    | 10      |

## Local

Stores your compose revisions locally inside the `.hcstate` directory next to your `helm-compose.yaml`.

| Option | Type   | Description                                                                                           | Default    |
| ------ | ------ | ----------------------------------------------------------------------------------------------------- | ---------- |
| `path` | string | The directory path to store your revisions (Relative to the directory you execute `helm compose` in). | `.hcstate` |

## Kubernetes

Stores your compose revisions similar to helm releases inside secrets in a kubernetes cluster namespace.

| Option        | Type   | Description                                        | Default         |
| ------------- | ------ | -------------------------------------------------- | --------------- |
| `namespace`   | string | The namespace to store your revisions in.          | default         |
| `kubeconfig`  | string | The path to your kubeconfig file                   | ~/.kube/config  |
| `kubecontext` | string | The context to use from your specified kubeconfig. | current-context |

## S3

Not yet implemented

## GCS

Not yet implemented
