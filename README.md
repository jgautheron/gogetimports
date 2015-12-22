# gogetimports

Get a JSON-formatted list of imports for a given Golang project.

### Get Started

    $ go get github.com/jgautheron/gogetimports
    $ gogetimports -path $GOPATH/src/github.com/Sirupsen/logrus

### Usage

```
Usage:

  gogetimports -path <directory>

Flags:

  -path                path to be scanned for imports
  -only-third-parties  return only third party imports
```