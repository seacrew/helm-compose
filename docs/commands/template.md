# helm compose template

Template all kubernetes resources for the releases specified in your `helm-compose.yaml` and print them out to stdout.

## Usage

The following command will print out all kubernetes resources that would be installed or upgraded on stdout.

```
helm compose template [releases...] [flags]
```

## Options

```
Flags:
  -h, --help   help for template

Global Flags:
  -f, --file string   Compose configuration file
```
