# Auto-Staging-Builder

[![Maintainability](https://api.codeclimate.com/v1/badges/b7d5203ef3e07f1538a9/maintainability)](https://codeclimate.com/github/auto-staging/builder/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/b7d5203ef3e07f1538a9/test_coverage)](https://codeclimate.com/github/auto-staging/builder/test_coverage)
[![GoDoc](https://godoc.org/github.com/auto-staging/builder?status.svg)](https://godoc.org/github.com/auto-staging/builder)
[![Go Report Card](https://goreportcard.com/badge/github.com/auto-staging/builder)](https://goreportcard.com/report/github.com/auto-staging/builder)
[![Build Status](https://travis-ci.com/auto-staging/builder.svg?branch=master)](https://travis-ci.com/auto-staging/builder)

> Builder creates the CodeBuild Jobs and the CloudWatchEvents rules for the Environments, it also starts the CodeBuild Job.

## Request Bodys for CodeBuild

### Create | Tower -> Builder

```json
{
  "operation": "CREATE",
  "repository": "my-app",
  "branch": "feat/test",
  "codeBuildRoleARN": "arn:aws:iam::123456789012:role/RepositoryCodeBuildRole",
  "infrastructureRepoUrl": "https://github.com/username/repository.git",
  "environmentVariables": [
    {
      "name": "TF_VAR_instance_type",
      "type": "PLAINTEXT",
      "value": "t2.micro"
    }
  ]
}
```

### After Create result | CodeBuild -> Builder

```json
{
  "operation": "RESULT_CREATE",
  "success": 1,
  "repository": "my-app",
  "branch": "feat/test",
}
```

### Update | Tower -> Builder

```json
{
  "operation": "UPDATE",
  "repository": "my-app",
  "branch": "feat/test",
  "codeBuildRoleARN": "arn:aws:iam::123456789012:role/RepositoryCodeBuildRole",
  "infrastructureRepoUrl": "https://github.com/username/repository.git",
  "environmentVariables": [
    {
      "name": "TF_VAR_instance_type",
      "type": "PLAINTEXT",
      "value": "t2.micro"
    }
  ]
}
```

### After Update result | CodeBuild -> Builder

```json
{
  "operation": "RESULT_UPDATE",
  "success": 1,
  "repository": "my-app",
  "branch": "feat/test",
}
```

### Delete | Tower -> Builder

```json
{
  "operation": "DELETE",
  "repository": "my-app",
  "branch": "feat/test"
}
```

### After Delete result | CodeBuild -> Builder

```json
{
  "operation": "RESULT_DELETE",
  "success": 1,
  "repository": "my-app",
  "branch": "feat/test",
}
```

## Request Bodys for CloudWatchEvents

### Create / Update Schedule Event | Tower -> Builder

```json
{
  "operation": "UPDATE_SCHEDULE",
  "repository": "my-app",
  "branch": "feat/test",
  "shutdownSchedules": [
    {
      "cron": "(0 12 * * ? *)"
    }
  ],
  "startupSchedules": [
    {
      "cron": "(0 11 * * ? *)"
    }
  ]
}
```

### Delete Schedule Event | Tower -> Builder

```json
{
  "operation": "DELETE_SCHEDULE",
  "repository": "my-app",
  "branch": "feat/test"
}
```

## Requirements

- Golang > 1.11

## Usage

### Run application

```bash
make run
```

### Build binary

```bash
make build
```

compiles to bin/auto-staging-builder

## License and Author

Author: Jan Ritter

License: MIT