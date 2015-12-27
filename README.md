# gogetimports

Get a JSON-formatted list of imports.

This tool will be useful if you'd like to get a bird view of the packages used by your application, or get statistics about third party libraries.

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

### License
MIT
