package model

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"

	"github.com/auto-staging/builder/helper"
	"github.com/auto-staging/builder/types"
)

// UpdateCloudWatchEvents removes all existing CloudWatchEvents rules for the Environment and creates new shutdown and startup schedules according
// to the Environment configuration. This function is called by the controller.
// If an error occurs the error gets logged and the returned.
func (CloudWatchEventsModel *CloudWatchEventsModel) UpdateCloudWatchEvents(event types.Event) error {
	// Adapt branch name to only contain allowed characters for CodeBuild name
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/UpdateCloudWatchEventsForEnvironment", "operation": "regex/compile"}, 0)
		return err
	}

	branchName := reg.ReplaceAllString(event.Branch, "-")

	// Startup schedules
	helper.Logger.Log(errors.New("Startup schedules: "+fmt.Sprint(event.StartupSchedules)), map[string]string{"module": "model/UpdateCloudWatchEvents", "operation": "startupSchedules"}, 4)

	err = CloudWatchEventsModel.removeRulesWithTarget(event.Repository, branchName, "start")
	if err != nil {
		return err
	}

	err = CloudWatchEventsModel.createRulesWithTarget(event.Repository, branchName, event.Branch, "start", event.StartupSchedules)
	if err != nil {
		return err
	}

	// Shutdown schedules
	helper.Logger.Log(errors.New("Shutdown schedules: "+fmt.Sprint(event.ShutdownSchedules)), map[string]string{"module": "model/UpdateCloudWatchEvents", "operation": "shutdownSchedules"}, 4)

	err = CloudWatchEventsModel.removeRulesWithTarget(event.Repository, branchName, "stop")
	if err != nil {
		return err
	}

	err = CloudWatchEventsModel.createRulesWithTarget(event.Repository, branchName, event.Branch, "stop", event.ShutdownSchedules)

	return err
}

func (CloudWatchEventsModel *CloudWatchEventsModel) removeRulesWithTarget(repository, branch, action string) error {
	helper.Logger.Log(errors.New("Removing "+action+" schedules for repo = "+repository+" and branch = "+branch), map[string]string{"module": "model/removeRuleWithTarget", "operation": "info/removeRules"}, 3)

	client := CloudWatchEventsModel.CloudWatchEventsAPI

	result, err := client.ListRules(nil)
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/removeRuleWithTarget", "operation": "aws/listRules"}, 0)
		return err
	}

	for i := range result.Rules {
		if strings.Contains(*result.Rules[i].Name, "as-"+action+"-"+repository+"-"+branch) {
			targetResult, err := client.ListTargetsByRule(&cloudwatchevents.ListTargetsByRuleInput{
				Rule: result.Rules[i].Name,
			})
			if err != nil {
				helper.Logger.Log(err, map[string]string{"module": "model/removeRulesWithTarget", "operation": "aws/listTargets"}, 0)
				return err
			}
			helper.Logger.Log(errors.New(fmt.Sprint(targetResult)), map[string]string{"module": "model/removeRulesWithTarget", "operation": "aws/listTargetsResult"}, 4)

			if len(targetResult.Targets) != 0 {
				helper.Logger.Log(errors.New("Deleting targets"), map[string]string{"module": "model/removeRulesWithTarget", "operation": "debug/deleteInfo"}, 4)
				var targetIds []*string

				for a := range targetResult.Targets {
					targetIds = append(targetIds, targetResult.Targets[a].Id)
				}
				deleteResult, err := client.RemoveTargets(&cloudwatchevents.RemoveTargetsInput{
					Ids:  targetIds,
					Rule: result.Rules[i].Name,
				})
				if err != nil {
					helper.Logger.Log(err, map[string]string{"module": "model/removeRulesWithTarget", "operation": "aws/removeTarget"}, 0)
					return err
				}
				helper.Logger.Log(errors.New(fmt.Sprint(deleteResult)), map[string]string{"module": "model/removeRulesWithTarget", "operation": "aws/removeTargetResult"}, 4)

			} else {
				helper.Logger.Log(errors.New("Skipping target delete, since there are no targets attached"), map[string]string{"module": "model/removeRulesWithTarget", "operation": "debug/targetDeleteInfo"}, 4)
			}

			helper.Logger.Log(errors.New("Deleting rule: "+*result.Rules[i].Name), map[string]string{"module": "model/removeRulesWithTarget", "operation": "info/deleteRule"}, 3)
			result, err := client.DeleteRule(&cloudwatchevents.DeleteRuleInput{
				Name: result.Rules[i].Name,
			})
			if err != nil {
				helper.Logger.Log(err, map[string]string{"module": "model/removeRulesWithTarget", "operation": "aws/deleteRule"}, 0)
				return err
			}
			helper.Logger.Log(errors.New(fmt.Sprint(result)), map[string]string{"module": "model/removeRulesWithTarget", "operation": "aws/removeTargetOutput"}, 4)
		}
	}

	return nil
}

// createRulesWithTarget creates a new CloudWatcheEvents rule with the Scheduler Lambda Funktion as target.
// If an error occurs the error gets logged and the returned.
func (CloudWatchEventsModel *CloudWatchEventsModel) createRulesWithTarget(repository, branch, branchRaw, action string, schedule []types.TimeSchedule) error {
	helper.Logger.Log(errors.New("Creating "+action+" schedules for repo = "+repository+" and branch = "+branchRaw), map[string]string{"module": "model/createRulesWithTarget", "operation": "info/createRules"}, 3)

	client := CloudWatchEventsModel.CloudWatchEventsAPI

	for i := range schedule {
		ruleName := "as-" + action + "-" + repository + "-" + branch + "-" + fmt.Sprint(rand.Intn(9999))

		helper.Logger.Log(errors.New("Adding rule with name: "+ruleName), map[string]string{"module": "model/createRulesWithTarget", "operation": "putRule"}, 4)
		output, err := client.PutRule(&cloudwatchevents.PutRuleInput{
			Description:        aws.String("Managed by auto-staging - cron" + schedule[i].Cron),
			Name:               aws.String(ruleName),
			ScheduleExpression: aws.String("cron" + schedule[i].Cron),
			RoleArn:            aws.String(os.Getenv("CLOUDWATCH_TO_LAMBDA_EXEC_ROLE")),
		})
		if err != nil {
			helper.Logger.Log(err, map[string]string{"module": "model/createRulesWithTarget", "operation": "aws/putRule"}, 0)
			return err
		}
		helper.Logger.Log(errors.New(fmt.Sprint(output)), map[string]string{"module": "model/createRulesWithTarget", "operation": "aws/putRuleOutput"}, 4)

		helper.Logger.Log(errors.New("Adding target to rule with name: "+ruleName), map[string]string{"module": "model/createRulesWithTarget", "operation": "putTarget"}, 4)
		target, err := client.PutTargets(&cloudwatchevents.PutTargetsInput{
			Targets: []*cloudwatchevents.Target{
				{
					Arn:   aws.String(os.Getenv("SCHEDULER_LAMBDA_ARN")),
					Id:    aws.String("scheduler" + fmt.Sprint(rand.Intn(9999))),
					Input: aws.String("{ \"repository\": \"" + repository + "\", \"branch\": \"" + branchRaw + "\", \"action\": \"" + action + "\" }"),
				},
			},
			Rule: aws.String(ruleName),
		})
		if err != nil {
			helper.Logger.Log(err, map[string]string{"module": "model/createRulesWithTarget", "operation": "aws/putTargets"}, 0)
			return err
		}
		helper.Logger.Log(errors.New(fmt.Sprint(target)), map[string]string{"module": "model/createRulesWithTarget", "operation": "aws/putTargetsOutput"}, 4)

	}
	return nil
}
