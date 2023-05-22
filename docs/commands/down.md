# helm compose down

Uninstall all releases defined in your `helm-compose.yaml`

## Usage

The following command will uninstall all releases of the previous revision if one exists. Otherwise the releases defined in your current `helm-compose.yaml` will be uninstalled.

```
helm compose down [flags]
```

## Options

```
Flags:
  -h, --help   help for down

Global Flags:
  -f, --file string   Compose configuration file
```
