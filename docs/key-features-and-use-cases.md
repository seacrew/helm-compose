# Key Features and Use Cases

The main idea behind `helm-compose` is to control / configure all helm related options as code by providing a compose file structure to configure everything you need to setup your helm based infrastructure.

## Repository handling

Configuration based installation of all necessary repositories you define as a dependency in your `helm-compose.yaml` before triggering the installation of your releases.

```yaml
apiVersion: 1.0

repositories:
  bitnami: https://charts.bitnami.com/bitnami
```

## Multi release handling

The main feature of `helm-compose` is the ability to define a multitude of releases inside a single file. `helm-compose` supports single kubernetes-cluster and multi-cluster configuration.

### Single cluster

Define as many releases as you would like for one or more namespaces.

```yaml
apiVersion: 1.0

releases:
  wordpress:
    chart: bitnami/wordpress
    chartVersion: 14.3.2
    namespace: homepage
  wordpress2:
    chart: bitnami/wordpress
    chartVersion: 15.2.22
    namespace: homepage
repositories:
  bitnami: https://charts.bitnami.com/bitnami
```

### Multi cluster

You can either use the `kubeconfig` options to point to a different path and use the `kubecontext` to select a specific context inside your kubeconfig.

```yaml
apiVersion: 1.0

releases:
  wordpress-dev:
    chart: bitnami/wordpress
    chartVersion: 14.3.2
    namespace: homepage
    kubeconfig: ~/.kube/dev
  wordpress-int:
    chart: bitnami/wordpress
    chartVersion: 15.2.22
    namespace: homepage
    kubeconfig: ~/.kube/int
repositories:
  bitnami: https://charts.bitnami.com/bitnami
```

### Environment variables

`helm-compose` is able to inject environment variables inside your values block to deal with secrets that shouldn't be committed to your source control.

Syntax: `${MY_ENV_VARIABLE}`.

```bash
export WORDPRESS_ADMIN_PASSWORD="xxx"
export MARIADB_ROOT_PASSWORD="xxx"
helm compose up
```

```yaml
apiVersion: 1.0

releases:
  wordpress:
    chart: bitnami/wordpress
    values:
      wordpressPassword: ${WORDPRESS_ADMIN_PASSWORD}
      mariadb.auth.rootPassword: ${MARIADB_ROOT_PASSWORD}
```

## Revision handling

Revisions are essentially snapshots of your current `helm-compose.yaml`. Every time you trigger `helm compose up` a new revision will be created and stored (By default the last 10 revisions are kept).

### Configuration

```yaml
apiVersion: 1.0

storage:
  name: wordpress
  type: local # default: local
  numberOfRevisions: 5 # default: 10
```

### Usage

You can list your revisions and get the content of your previous revisions via the [`helm compose list`](commands/list.md) and [`helm compose get`](commands/get.md) commands.

### Rollback

Just select a revision you want to go back to. You can check the content with the `helm compose get` command and parse it directly into any `helm compose` command like so:

```bash
$ helm compose list
| Date             | Revision |
| ---------------- | -------- |
| 2023-05-24 23:56 |       12 |
| 2023-05-24 23:56 |       13 |
| 2023-05-24 23:57 |       14 |
| 2023-05-24 23:57 |       15 |
| 2023-05-24 23:57 |       16 |

# select revision 15 and use the pipe | operator to parse the content back into compose up with -f -
helm compose get 15 | helm compose up -f -
```
