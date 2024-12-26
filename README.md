# Nushell Plugin

[Nushell](https://www.nushell.sh/)
[Plugin](https://www.nushell.sh/contributor-book/plugins.html) 
written in [Go](https://go.dev/) using 
[nu-plugin package](https://github.com/ainvaltin/nu-plugin).

## Implements commands

- `to plist` and `from plist` - convert to and from [Property List](https://en.wikipedia.org/wiki/Property_list) format;
- `encode base85` and `decode base85` - encode and decode [ascii85 / base85](https://en.wikipedia.org/wiki/Ascii85) encoded data;
- `encode base58` and `decode base58`;

Note that Nu has introduced `plist` support in version 0.97 of it's
["native" formats plugin](https://www.nushell.sh/commands/categories/formats.html).

## Installation

Latest version is for Nushell version `0.101.0`.

To install it you need to have [Go installed](https://go.dev/dl/), then run
```sh
go install github.com/ainvaltin/nu_plugin_plist@latest
```
This creates the `nu_plugin_plist` binary in your `GOBIN` directory:

> Executables are installed in the directory named by the GOBIN environment
variable, which defaults to $GOPATH/bin or $HOME/go/bin if the GOPATH
environment variable is not set.

Locate the binary and follow instructions on 
[Downloading and installing a plugin](https://www.nushell.sh/book/plugins.html#downloading-and-installing-a-plugin)
page on how to register `nu_plugin_plist` as a plugin.
