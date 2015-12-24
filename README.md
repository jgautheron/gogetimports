# gogetimports

Get a JSON-formatted list of imports.

### Get Started

    $ go get github.com/jgautheron/gogetimports
    $ gogetimports ./...

### Usage

```
Usage:

  gogetimports ARGS <directory>

Flags:

  -only-third-parties  return only third party imports

Examples:

  gogetimports ./...
  gogetimports -only-third-parties $GOPATH/src/github.com/cockroachdb/cockroach
```