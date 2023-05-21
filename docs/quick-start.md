# Quick Start

## Installation

Install a specific version (recommended). Click [here](https://github.com/seacrew/helm-compose/releases/latest) for the latest. version.

```
helm plugin install https://github.com/seacrew/helm-compose --version 1.0.0-beta.1
```

Install latest unstable version from main branch.

```
helm plugin install https://github.com/seacrew/helm-compose
```

## How to use helm compose

Helm Compose makes it easy to define a Compose file containing a list of Releases and necessary Repositories for the charts you use.

Install your releases:

```
helm compose up -f helm-compose.yaml
```

Uninstall your releases

```
helm compose down -f helm-compose.yaml
```

A Helm Compose file looks something like this:

```yaml
apiVersion: 1.0

storage:
  name: mycompose
  type: local # default
  path: .hcstate # default

releases:
  wordpress:
    chart: bitnami/wordpress
    chartVersion: 14.3.2
  wordpress2:
    chart: bitnami/wordpress
    chartVersion: 15.2.22
    namespace: homepage
    createNamespace: true
  postgres:
    chart: bitnami/postgresql
    chartVersion: 12.1.9
    namespace: database
    createNamespace: true

repositories:
  bitnami: https://charts.bitnami.com/bitnami
```

Check out the [helm compose examples](https://github.com/seacrew/helm-compose/tree/main/examples).
