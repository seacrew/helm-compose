![helm-compose-banner](https://user-images.githubusercontent.com/18513179/212496531-1d166236-ed88-411d-8403-ad1f94d28846.png)
# Helm Compose

# IMPORTANT NOTICE: THIS IS A WORK IN PROGRESS
As of yet, this plugin should __NOT__ be treated as a stable tool and shouldn't be used in an production environment.

## About
Helm Compose is a tool for managing multiple releases for one or many different Helm charts. It is an extension of the package manager idea behind Helm and is heavily inspired by Docker Compose.

# Installation
Install a specific version (recommended). Click [here](https://github.com/seacrew/helm-compose/releases/latest) for the latest. version.
```
helm plugin install https://github.com/seacrew/helm-compose --version 1.0.0-alpha.2
```

Install latest unstable version from main branch.
```
helm plugin install https://github.com/seacrew/helm-compose
```

# Quick Start
Helm Compose makes it easy to define a Compose file containing a list of Releases and necessary Repositories for the charts you use.

A Compose file looks something like this:

```yaml
apiVersion: 1.0

compose:
  name: mycompose
  state: local

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

Install your releases: 
```bash
$ helm compose up -f helm-compose.yaml
```

Uninstall your releases
```bash
$ helm compose down -f helm-compose.yaml
```