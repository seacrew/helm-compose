# Storage Providers

Following options are applied regardless of the selected provider.

| Option            | Type   | Description | Required | Default |
| ----------------- | ------ | ----------- | -------- | ------- |
| name              | string |             | true     |         |
| type              | string |             | false    | local   |
| numberOfRevisions | int    |             | false    | 10      |

## Local

Stores your compose revision locally inside the `.hcstate` directory next to your `helm-compose.yaml`.

| Option | Type   | Description                                                                                           | Default    |
| ------ | ------ | ----------------------------------------------------------------------------------------------------- | ---------- |
| `path` | string | The directory path to store your revisions (Relative to the directory you execute `helm compose` in). | `.hcstate` |

## Kubernetes

Not yet implemented

| Option      | Type   | Description                               | Default |
| ----------- | ------ | ----------------------------------------- | ------- |
| `namespace` | string | The namespace to store your revisions in. | default |

## S3

Not yet implemented
