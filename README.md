## SEctl

SEctl is a binary used to query SELinux and display information related to your current configuration / policy.

SEctl is _not _a rewrite of the current SELinux tools. It's just a more simple way to query SELinux information from one single binary.

### Example

```
sectl status | jq .mode
"permissive"
```

### Build

```
make
```

### Install

Copy `sectl` in your path.

### Getting help

```
$ sectl help
sectl is a tool to query SELinux

Usage:
  sectl [command]

Available Commands:
  help        Help about any command
  status      display the status of SELinux

Flags:
  -h, --help            help for sectl
      --output string   output format (default "json")

Use "sectl [command] --help" for more information about a command.
```

