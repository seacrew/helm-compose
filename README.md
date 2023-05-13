![helm-compose-banner](https://user-images.githubusercontent.com/18513179/212496531-1d166236-ed88-411d-8403-ad1f94d28846.png)

# Helm Compose

[![Build Status](https://github.com/seacrew/helm-compose/actions/workflows/build.yml/badge.svg)](https://github.com/seacrew/helm-compose/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/seacrew/helm-compose)](https://goreportcard.com/report/github.com/seacrew/helm-compose)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/seacrew/helm-compose)

# IMPORTANT NOTICE: THIS IS A WORK IN PROGRESS

As of yet, this plugin should **NOT** be treated as a stable tool and shouldn't be used in an production environment.

## About

Helm Compose is a tool for managing multiple releases for one or many different Helm charts. It is an extension of the package manager idea behind Helm and is heavily inspired by Docker Compose.

## Installation

Install a specific version (recommended). Click [here](https://github.com/seacrew/helm-compose/releases/latest) for the latest. version.

```
helm plugin install https://github.com/seacrew/helm-compose --version 1.0.0-alpha.3
```

Install latest unstable version from main branch.

```
helm plugin install https://github.com/seacrew/helm-compose
```

## Quick Start Guide

Helm Compose makes it easy to define a Compose file containing a list of Releases and necessary Repositories for the charts you use.

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
apiVersion: 1.0

state:
  name: mycompose
  storage: local

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

## Compose file reference

The compose file is a [YAML](https://yaml.org/) file for defining your Helm releases and necessary Helm repositories. This allows you manage multiple releases. Scenarios in which this might be useful:

- Multiple deployments of the same helm chart but with different versions (in the same or different namespaces)
- Multiple interdependent helm chart deployments (in the same or different namespaces)
- Multiple deployments in multiple k8s clusters

### `state`

To keep track of the changes in your compose file. A snapshot / state of your compose file is stored every time you make an update. Similar to how helm stores single release information as a secret inside the same namespace as the release. With helm compose you have two options:

`local` (default): Stores the 10 lastest states inside the .hcstate directory

```yaml
state:
  name: mycompose
  storage: local
```

`kubernetes`: Similar to helm itself this option stores the 10 latest states as secrets inside a specified namespace

```yaml
state:
  name: mycompose
  storage: kubernetes
  namespace: states
```

### `releases`

You can define as many releases as you want to treat them as a single entity. All fields are optional except for the chart.

```yaml
releases:
  wordpress:
    chart: your-chart
    chartVersion: 1.0.0
    namespace: your-namespace
    createNamespace: false
    kubeconfig: # Path to your custom Kubeconfig
    kubecontext: your-kube-context
    valuefiles: [] # list of value files to be used
    values: # your custom values
      key: value
```

### `repositories`

You can define as many repositories as you want and need for your charts to be available.

```yaml
repositories:
  bitnami: https://charts.bitnami.com/bitnami
```
