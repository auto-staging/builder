package model

import (
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"gitlab.com/auto-staging/builder/types"
	yaml "gopkg.in/yaml.v2"
)

func CreateCodeBuildJob(event types.Event) error {
	envVars := []*codebuild.EnvironmentVariable{}
	for key, value := range event.EnvironmentVariables {
		envVars = append(envVars, &codebuild.EnvironmentVariable{
			Name:  aws.String(key),
			Type:  aws.String("PLAINTEXT"),
			Value: aws.String(value),
		})
	}

	buildspec := types.Buildspec{
		Version: "0.2",
		Phases: types.Phases{
			Build: types.Build{
				Commands: []string{
					"terraform --version",
					"aws lambda invoke --function-name auto-staging-builder --invocation-type Event --payload '{ \"operation\": \"SUCCESS_CREATE\", \"repository\": \"" + event.Repository + "\", \"branch\": \"" + event.Branch + "\" }'  /dev/null",
				},
			},
		},
	}

	// Adapt branch name to only contain allowed characters for CodeBuild name
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	branchName := reg.ReplaceAllString(event.Branch, "-")

	res, err := yaml.Marshal(buildspec)

	log.Println(err)
	log.Println(string(res))

	createInput := codebuild.CreateProjectInput{
		Name:        aws.String("auto-staging-" + event.Repository + "-" + branchName),
		Description: aws.String("Managed by auto-staging"),
		ServiceRole: aws.String("arn:aws:iam::171842373341:role/auto-staging-builder-codebuild-exec-role"),
		Environment: &codebuild.ProjectEnvironment{
			ComputeType:          aws.String("BUILD_GENERAL1_SMALL"),
			Image:                aws.String("janrtr/auto-staging-build"),
			Type:                 aws.String("LINUX_CONTAINER"),
			EnvironmentVariables: envVars,
		},
		Source: &codebuild.ProjectSource{
			Type:      aws.String("GITHUB"),
			Location:  aws.String(event.RepositoryURL),
			Buildspec: aws.String(string(res)),
		},
		Artifacts: &codebuild.ProjectArtifacts{
			Type: aws.String("NO_ARTIFACTS"),
		},
	}

	client := getCodeBuildClient()
	_, err = client.CreateProject(&createInput)
	if err != nil {
		log.Println(err)
	}

	return err
}
