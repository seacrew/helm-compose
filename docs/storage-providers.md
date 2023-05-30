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

Stores your compose revisions inside a s3 bucket.

| Option           | Type   | Description                                                                                                | Required | Default                                                                                                 |
| ---------------- | ------ | ---------------------------------------------------------------------------------------------------------- | -------- | ------------------------------------------------------------------------------------------------------- |
| s3bucket         | string | Specify the bucket name to upload to and download from.                                                    | true     |                                                                                                         |
| s3prefix         | string | Specify the object prefix (directory) for the revisions to be uploaded to and downloaded from.             | false    |                                                                                                         |
| s3region         | string | Set a custom S3 region.                                                                                    | false    | By default the region will be read from AWS configuration files or the AWS_REGION environment variable. |
| s3endpoint       | string | Set a custom S3 endpoint / host url.                                                                       | false    | Default AWS S3 service endpoint based on the region.                                                    |
| s3insecure       | bool   | Disable the verification of the servers certificate chain and hostname.                                    | false    | false                                                                                                   |
| s3disableSSL     | bool   | Disable the usage of SSL / https.                                                                          | false    | false                                                                                                   |
| s3forcePathStyle | bool   | Enforce to use path style. Especially useful for none AWS S3 provider which often only support path style. | false    | false                                                                                                   |

## GCS

Not yet implemented
