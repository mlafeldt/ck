# ck - The ConvertKit Tool

Access the [ConvertKit API](http://help.convertkit.com/article/33-api-documentation-v3) from the command line.

## Installation

You can build the tool from source, provided you have Go installed:

```bash
go get -u github.com/mlafeldt/ck
```

## Usage

List all confirmed subscribers for the given Convertkit account:

```bash
ck [-csv]
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
