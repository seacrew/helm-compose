# Storage Providers

Following options are applicable regardless of the selected provider.

| Option            | Type   | Description                                                                                 | Required | Default | ApiVersion |
| ----------------- | ------ | ------------------------------------------------------------------------------------------- | -------- | ------- | ---------- |
| name              | string | Name to be used to store revisions with a storage provider.                                 | true     |         | 1.0        |
| type              | string | Type / name of the storage provider you want to use. By default local files will be stored. | false    | local   | 1.0        |
| numberOfRevisions | int    | Number of revisions to be stored and to be able to rollback to.                             | false    | 10      | 1.0        |

## Local

Stores your compose revisions locally inside the `.hcstate` directory next to your `helm-compose.yaml`.

| Option | Type   | Description                                                                                           | Required | Default  | ApiVersion |
| ------ | ------ | ----------------------------------------------------------------------------------------------------- | -------- | -------- | ---------- |
| path   | string | The directory path to store your revisions (Relative to the directory you execute `helm compose` in). | false    | .hcstate | 1.0        |

## Kubernetes

Stores your compose revisions similar to helm releases inside secrets in a kubernetes cluster namespace.

| Option      | Type   | Description                                        | Required | Default         | ApiVersion |
| ----------- | ------ | -------------------------------------------------- | -------- | --------------- | ---------- |
| namespace   | string | The namespace to store your revisions in.          | false    | default         | 1.0        |
| kubeconfig  | string | The path to your kubeconfig file                   | false    | ~/.kube/config  | 1.0        |
| kubecontext | string | The context to use from your specified kubeconfig. | false    | current-context | 1.0        |

## S3

Stores your compose revisions inside a s3 bucket. You will need to set your AWS credentials (access and secret key) via [environment variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html) or the `~/.aws/config` file. [Official AWS documentation](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html).


| Option           | Type   | Description                                                                                                | Required | Default                                                                      | ApiVersion |
| ---------------- | ------ | ---------------------------------------------------------------------------------------------------------- | -------- | ---------------------------------------------------------------------------- | ---------- |
| s3bucket         | string | Specify the bucket name to upload to and download from.                                                    | true     |                                                                              | 1.0        |
| s3prefix         | string | Specify the object prefix (directory) for the revisions to be uploaded to and downloaded from.             | false    | (root path)                                                                  | 1.0        |
| s3region         | string | Set a custom S3 region.                                                                                    | false    | By default the region will be read from the AWS_REGION environment variable. | 1.0        |
| s3endpoint       | string | Set a custom S3 endpoint / host url.                                                                       | false    | Default AWS S3 service endpoint.                                             | 1.0        |
| s3insecure       | bool   | Disable the verification of the servers certificate chain and hostname.                                    | false    | false                                                                        | 1.0        |
| s3disableSSL     | bool   | Disable the usage of SSL / https.                                                                          | false    | false                                                                        | 1.0        |
| s3forcePathStyle | bool   | Enforce to use path style. Especially useful for none AWS S3 provider which often only support path style. | false    | false                                                                        | 1.0        |

## GCS

Not yet implemented
