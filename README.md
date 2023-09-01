![helm-compose-banner](https://user-images.githubusercontent.com/18513179/240495789-e76890d3-f0f9-48b9-9d18-89e53effe65b.png)

[![Build Status](https://github.com/seacrew/helm-compose/actions/workflows/build.yaml/badge.svg)](https://github.com/seacrew/helm-compose/actions/workflows/build.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/seacrew/helm-compose)](https://goreportcard.com/report/github.com/seacrew/helm-compose)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=seacrew_helm-compose&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=seacrew_helm-compose)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=seacrew_helm-compose&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=seacrew_helm-compose)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/seacrew/helm-compose)](https://github.com/seacrew/helm-compose/releases/latest)

Helm Compose is a tool for managing multiple releases of one or many different Helm charts. It is heavily inspired by Docker Compose and is an extension of the package manager idea behind Helm itself.

## Installation

It is requirement to use helm v3.10.0+. 

Install a specific version of helm compose (recommended). Click [here](https://github.com/seacrew/helm-compose/releases/latest) for the latest version.

```
helm plugin install https://github.com/seacrew/helm-compose --version 1.2.0
```

Install the latest version.

```
helm plugin install https://github.com/seacrew/helm-compose
```

## Quick Start Guide

Helm Compose makes it easy to define a list of Releases and all necessary Repositories for the charts you use in a single compose file.

Install your releases:

```bash
$ helm compose up -f helm-compose.yaml
```

Uninstall your releases

```bash
$ helm compose down -f helm-compose.yaml
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

Check out the [examples](https://github.com/seacrew/helm-compose/tree/main/examples) directory.

## Documentation

Checkout the complete [documentation.](https://seacrew.github.io/helm-compose/)
