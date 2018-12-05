# Auto-Staging-Builder

## Request Bodys

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

## Usage

#### Install dependencies

```bash
make prepare
```

#### Run application

```bash
make run
```

#### Build binary

```bash
make build
```

compiles to bin folder

## License and Author

Author: Jan Ritter

License: MIT