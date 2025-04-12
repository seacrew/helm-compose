# Quick Start

## Installation

Install a specific version (recommended). Click [here](https://github.com/seacrew/helm-compose/releases/latest) for the latest version.

```
helm plugin install https://github.com/seacrew/helm-compose --version 1.4.0
```

Install latest version.

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
apiVersion: 1.1

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

All `helm-compose` commands accept the `-f` flag to pass your compose file location. Otherwise `helm-compose` will automatically look for a list of file names inside your current directory:

- helm-compose.yaml
- helm-compose.yml
- helmcompose.yaml
- helm-compose.yml
- helmcompose.yaml
- helmcompose.yml
- helmcompose
- compose.yaml
- compose.yml

Check out the [helm compose examples](https://github.com/seacrew/helm-compose/tree/main/examples).
