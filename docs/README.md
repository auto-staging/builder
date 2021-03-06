# Auto-Staging-Builder

[![Maintainability](https://api.codeclimate.com/v1/badges/b7d5203ef3e07f1538a9/maintainability)](https://codeclimate.com/github/auto-staging/builder/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/b7d5203ef3e07f1538a9/test_coverage)](https://codeclimate.com/github/auto-staging/builder/test_coverage)
[![GoDoc](https://godoc.org/github.com/auto-staging/builder?status.svg)](https://godoc.org/github.com/auto-staging/builder)
[![Go Report Card](https://goreportcard.com/badge/github.com/auto-staging/builder)](https://goreportcard.com/report/github.com/auto-staging/builder)
[![Build Status](https://travis-ci.com/auto-staging/builder.svg?branch=master)](https://travis-ci.com/auto-staging/builder)

> Builder creates the CodeBuild Jobs and the CloudWatchEvents rules for the Auto Staging Environments, it also starts the CodeBuild Job.

## Requirements

- Golang > 1.11

## Usage

### Run tests

```bash
make tests
```

### Build binary

```bash
make build
```

compiles to bin/auto-staging-builder

## Request bodies

### [Request bodies for CodeBuild](REQUEST-BODIES-CODEBUILD.md)

### [Request bodies for CloudWatch](REQUEST-BODIES-CLOUDWATCH.md)

## License and Author

Author: Jan Ritter

License: MIT