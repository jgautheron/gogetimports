# gogetimports

Get a JSON-formatted map of imports per file.

This tool will be useful if you'd like to get a bird view of the packages used by your application, or get statistics about third party libraries used.

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
  -pretty              output JSON with proper indentation

Examples:

  gogetimports ./...
  gogetimports -only-third-parties $GOPATH/src/github.com/cockroachdb/cockroach
  gogetimports -ignore "jgautheron\/gocha" -list $GOPATH/src/github.com/jgautheron/gocha/...
  gogetimports -pretty .
```

### License
MIT
