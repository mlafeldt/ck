# ck - The ConvertKit Tool

[![Build Status](https://travis-ci.org/mlafeldt/ck.svg?branch=master)](https://travis-ci.org/mlafeldt/ck)
[![GoDoc](https://godoc.org/github.com/mlafeldt/ck/convertkit?status.svg)](https://godoc.org/github.com/mlafeldt/ck/convertkit)

Access the [ConvertKit API](http://help.convertkit.com/article/33-api-documentation-v3) from the command line.

## Installation

If you're on Mac OS X, the easiest way to get the `ck` command-line tool is via Homebrew:

```bash
brew tap mlafeldt/formulas
brew install cktool
```

You can also build the tool from source, provided you have Go installed:

```bash
go get -u github.com/mlafeldt/ck
```

## Usage

```
Usage:
  ck [command]

Available Commands:
  help        Help about any command
  subscribers List all confirmed subscribers
  version     Show program version

Flags:
  -h, --help   help for ck

Use "ck [command] --help" for more information about a command.
```

The tool understands these environment variables:

* `CONVERTKIT_ENDPOINT`
* `CONVERTKIT_API_KEY`
* `CONVERTKIT_API_SECRET`

## Go library

In addition to the CLI tool, the project also provides the `convertkit` Go library for use in other Go projects. To install it from source:

```bash
go get -u github.com/mlafeldt/ck/convertkit
```

For usage and examples, see the [Godoc documentation](https://godoc.org/github.com/mlafeldt/ck/convertkit).

## Author

This project is being developed by [Mathias Lafeldt](https://twitter.com/mlafeldt).
