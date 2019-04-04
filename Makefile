
prepare:
	dep ensure -v

build: prepare
	GOOS=linux go build -o ./bin/auto-staging-builder -v -ldflags "-X github.com/auto-staging/builder/helper.commitHash=`git rev-parse HEAD` -X github.com/auto-staging/builder/helper.buildTime=`date -u +"%Y-%m-%dT%H:%M:%SZ"` -X github.com/auto-staging/builder/helper.branch=`git rev-parse --abbrev-ref HEAD` -X github.com/auto-staging/builder/helper.version=`git describe --abbrev=0 --tags`"

tests:
	go test ./... -v -cover

run:
	go run main.go
