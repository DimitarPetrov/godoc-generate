# godoc-generate

[![Go Report Card](https://goreportcard.com/badge/github.com/DimitarPetrov/godoc-generate)](https://goreportcard.com/report/github.com/DimitarPetrov/godoc-generate)

## Overview

`godoc-generate` is a simple command line tool that generates default godoc comments on all **exported** `types`, `functions`, `consts` and `vars` in the current working directory and recursively for all subdirectories.

The godoc comments looks like this:

```
// %s missing godoc.
```

Where `%s` is the name of the type/func/const/var.

## Installation

#### Installing from Source
```
go install github.com/DimitarPetrov/godoc-generate@latest
```

## Demonstration

Let's say you have a simple `Multiply` function without `godoc`:

```go
func Multiply(a,b int) int {
	return a * b
}
```

It is exported, therefore it is part of the package's interface. It is ideomatic to add godoc on everything exported in your package.

If you run `godoc-genenrate` the code will be rewritten the following way:

```go
// Multiply missing godoc.
func Multiply(a, b int) int {
	return a * b
}
```

This way you are safe to add a linter enforcing godoc and migrate all legacy code gradually.
