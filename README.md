# Replacer

[![Go Report Card](https://goreportcard.com/badge/github.com/weastur/replacer)](https://goreportcard.com/report/github.com/weastur/replacer)
[![codecov](https://codecov.io/gh/weastur/replacer/graph/badge.svg?token=QANQ7BIQY9)](https://codecov.io/gh/weastur/replacer)
[![test](https://github.com/weastur/replacer/actions/workflows/test.yaml/badge.svg)](https://github.com/weastur/replacer/actions/workflows/test.yaml)
[![lint](https://github.com/weastur/replacer/actions/workflows/lint.yaml/badge.svg)](https://github.com/weastur/replacer/actions/workflows/lint.yaml)
[![gitlint](https://github.com/weastur/replacer/actions/workflows/gitlint.yaml/badge.svg)](https://github.com/weastur/replacer/actions/workflows/gitlint.yaml)
[![pre-commit.ci status](https://results.pre-commit.ci/badge/github/weastur/replacer/main.svg)](https://results.pre-commit.ci/latest/github/weastur/replacer/main)</br>
![GitHub Release](https://img.shields.io/github/v/release/weastur/replacer)
![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/weastur/replacer/total)
![GitHub commits since latest release](https://img.shields.io/github/commits-since/weastur/replacer/latest)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/weastur/replacer)
![GitHub License](https://img.shields.io/github/license/weastur/replacer)

**Replacer** is a Go code generator that applies regex-based transformations to your source files.
It is designed to work with Go's `//go:generate` directive, allowing you to automate repetitive code modifications.

## Why?

I often need to repeat the same set of comments in my Go code.
For example, I use [swaggo](https://github.com/swaggo/swag) to generate Swagger documentation for an API.
My API endpoints always return the same headers, like `X-API-Version` or `X-Request-ID`.
Unfortunately, there's no way to define common headers for all endpoints yet.
Instead of manually copying and pasting these headers every time, I created `replacer` to automate this process.
Now, I can simply write:

```go
// @COMMON-HEADERS
```

run `go generate`, which is a part of my build pipeline, and get:

```go
// @Header       all {string} X-Request-ID "UUID of the request"
// @Header       all {string} X-API-Version "API version, e.g. v1alpha"
// @Header       all {int} X-Ratelimit-Limit "Rate limit value"
// @Header       all {int} X-Ratelimit-Remaining "Rate limit remaining"
// @Header       all {int} X-Ratelimit-Reset "Rate limit reset interval in seconds"
```

## Features

- Full-featured regex ([docs](https://pkg.go.dev/regexp))
- Automatically searches for a configuration file (`.replacer.yml` or `.replacer.yaml`)
    in the current directory and parent directories.
- Stops searching at the root directory or when a `go.mod` file is encountered.

## Installation

To install `replacer`, run:

```bash
go install github.com/weastur/replacer/cmd/replacer@latest
```

Make sure that `$GOPATH/bin` is in your `$PATH`.

Of course, you can also download the binary from the [releases page](https://github.com/weastur/replacer/releases)

## Usage

1. Create a configuration file (`.replacer.yml` or `.replacer.yaml`) in the root of your project
    or in the directory where you want to run `replacer`.

    ```yaml
    rules:
    - regex: '(?m)^// MY RULE$'
        repl: |-
        // MY NEW AWESOME
        // MULTI-LINE REPLACEMENT
    - regex: '(?m)^// MY 2d RULE$'
        repl: |-
        // ANOTHER REPLACEMENT
    ```

1. Add a `//go:generate` directive to your source file.

    ```go
    //go:generate replacer
    ```

1. Run `go generate` to apply the transformations.

    ```bash
    go generate ./...
    ```

    Pay attention that generate command will only work if you have the `replacer` binary in your `$PATH`.

    Also, the `go build` command doesn't run `go generate` automatically.
    You need to run it manually. Refer to `go help generate` for more information.

1. Optional Flags:

    - `-config` - specify the path to the configuration file.

        Example: `//go:generate replacer -config my-replacer-config.yml`

### Configuration file search

If the `-config` flag is not provided, `replacer` will search for a configuration file in the following order:

1. The current directory.
1. Parent directories, moving up one level at a time.
1. Stops searching when it reaches the root directory or encounters a go.mod file.

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md) for information on how to contribute to `replacer`.

### tl;dr

```bash
# fork/clone the repository
# install go, direnv, pre-commit hooks
# put some config in .replacer.yml
make build
make test
GOFILE=my-test-file.go replacer
# commit/push/PR
```

### Project structure

- `cmd/replacer` - the main command.
- `internal/config` - handles configuration file parsing and validation.
- `internal/replacer` - applies regex-based transformations to source files.

## Security

Refer to the [SECURITY.md](SECURITY.md) file for more information.

## License

Mozilla Public License 2.0

Refer to the [LICENSE](LICENSE) file for more information.
