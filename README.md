# gogetimports

Get a JSON-formatted map of imports per file.

### Get Started

    $ go get github.com/jgautheron/gogetimports
    $ gogetimports ./...

### Usage

```
Usage:

  gogetimports ARGS <directory>

Flags:

  -only-third-parties  return only third party imports
  -list                return a list instead of a map
  -ignore              ignore imports matching the given regular expression

Examples:

  gogetimports ./...
  gogetimports -only-third-parties $GOPATH/src/github.com/cockroachdb/cockroach
  gogetimports -ignore "jgautheron" -list $GOPATH/src/github.com/jgautheron/gocha/...
```

### License
MIT