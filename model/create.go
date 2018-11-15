package model

import (
	"log"

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
					"apt-get update && apt-get install golang -y",
				},
			},
		},
	}

	res, err := yaml.Marshal(buildspec)

	log.Println(err)
	log.Println(string(res))

	createInput := codebuild.CreateProjectInput{
		Name:        aws.String("auto-staging-" + event.Repository + "-" + event.Branch),
		Description: aws.String("Managed by auto-staging"),
		ServiceRole: aws.String("arn:aws:iam::171842373341:role/auto-staging-builder-codebuild-exec-role"),
		Environment: &codebuild.ProjectEnvironment{
			ComputeType:          aws.String("BUILD_GENERAL1_SMALL"),
			Image:                aws.String("ubuntu:18.04"),
			Type:                 aws.String("LINUX_CONTAINER"),
			EnvironmentVariables: envVars,
		},
		Source: &codebuild.ProjectSource{
			Type:      aws.String("GITHUB"),
			Location:  aws.String("https://github.com/janritter/kvb-api.git"),
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
