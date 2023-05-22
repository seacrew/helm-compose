# helm compose up

Install all releases and repositories defined in your `helm-compose.yaml`

## Usage

The following command will install all releases defined in your `helm-compose.yaml` and will compare it to the latest previous revision and uninstall all releases that have been removed since then.

```
helm compose up [flags]
```

## Options

```
Flags:
  -h, --help   help for up

Global Flags:
  -f, --file string   Compose configuration file
```
