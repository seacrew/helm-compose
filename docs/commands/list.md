# helm compose list

List previous revisions of your `helm-compose.yaml`

## Usage

```
helm compose list [flags]
```

## Options

```
Flags:
  -h, --help   help for list

Global Flags:
  -f, --file string   Compose configuration file
```

## Example

```
$ helm compose list

| Date             | Revision |
| ---------------- | -------- |
| 2023-05-22 10:24 |        1 |
| 2023-05-22 11:19 |        2 |
```
