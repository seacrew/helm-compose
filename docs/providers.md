# State Storage Providers

Following options are applied regardless of the selected provider.

| Option           | Type         | Description                           | Default |
| ---------------- | ------------ | ------------------------------------- | ------- |
| `name`           | string       | The name of your state files.         |         |
| `type`           | ProviderType | The name of a state storage provider. | `local` |
| `numberOfStates` | uint         | The number of states to be stored.    | 10      |

## Local

Stores your compose state locally inside the `.hcstate` directory next to your `helm-compose.yaml`.

| Option | Type   | Description                                                                                        | Default    |
| ------ | ------ | -------------------------------------------------------------------------------------------------- | ---------- |
| `path` | string | The directory path to store your states (Relative to the directory you execute `helm compose` in). | `.hcstate` |

## Kubernetes

Not yet implemented

| Option      | Type   | Description                            | Default |
| ----------- | ------ | -------------------------------------- | ------- |
| `namespace` | string | The namespace to store your states in. | default |

## S3

Not yet implemented
