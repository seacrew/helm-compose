# Storage Providers

Following options are applicable regardless of the selected provider.

| Option            | Type   | Description                                                                                 | Required | Default |
| ----------------- | ------ | ------------------------------------------------------------------------------------------- | -------- | ------- |
| name              | string | Name to be used to store revisions with a storage provider.                                 | true     |         |
| type              | string | Type / name of the storage provider you want to use. By default local files will be stored. | false    | local   |
| numberOfRevisions | int    | Number of revisions to be stored and to be able to rollback to.                             | false    | 10      |

## Local

Stores your compose revisions locally inside the `.hcstate` directory next to your `helm-compose.yaml`.

| Option | Type   | Description                                                                                           | Required | Default  |
| ------ | ------ | ----------------------------------------------------------------------------------------------------- | -------- | -------- |
| path   | string | The directory path to store your revisions (Relative to the directory you execute `helm compose` in). | false    | .hcstate |

## Kubernetes

Stores your compose revisions similar to helm releases inside secrets in a kubernetes cluster namespace.

| Option      | Type   | Description                                        | Required | Default         |
| ----------- | ------ | -------------------------------------------------- | -------- | --------------- |
| namespace   | string | The namespace to store your revisions in.          | false    | default         |
| kubeconfig  | string | The path to your kubeconfig file                   | false    | ~/.kube/config  |
| kubecontext | string | The context to use from your specified kubeconfig. | false    | current-context |

## S3

Not yet implemented

## GCS

Not yet implemented
