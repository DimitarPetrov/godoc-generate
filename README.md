# godoc-generate

[![Build Status](https://travis-ci.org/DimitarPetrov/godoc-generate.svg?branch=master)](https://travis-ci.org/DimitarPetrov/godoc-generate)
[![Coverage Status](https://coveralls.io/repos/github/DimitarPetrov/godoc-generate/badge.svg?branch=master)](https://coveralls.io/github/DimitarPetrov/godoc-generate?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/DimitarPetrov/godoc-generate)](https://goreportcard.com/report/github.com/DimitarPetrov/godoc-generate)

## Overview

`godoc-generate` is a simple command line tool that generates default godoc comments on all **exported** `types`, `functions`, `consts` and `vars` in the current working directory.

The godoc comments looks like this, where `%s` is the name of the type/func/const/var:

```
// %s missing godoc
```

## Installation

#### Installing from Source
```
go install github.com/DimitarPetrov/godoc-generate@latest
```
