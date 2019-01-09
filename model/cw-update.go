package model

import (
	"fmt"
	"log"
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
func UpdateCloudWatchEvents(event types.Event) error {
	// Adapt branch name to only contain allowed characters for CodeBuild name
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/UpdateCloudWatchEventsForEnvironment", "operation": "regex/compile"}, 0)
		return err
	}

	branchName := reg.ReplaceAllString(event.Branch, "-")

	// Startup schedules
	fmt.Println("Startup schedules - Target")
	fmt.Println(event.StartupSchedules)

	err = removeRulesWithTarget(event.Repository, branchName, "start")
	if err != nil {
		return err
	}

	err = createRulesWithTarget(event.Repository, branchName, event.Branch, "start", event.StartupSchedules)
	if err != nil {
		return err
	}

	// Shutdown schedules
	fmt.Println("Shutdown schedules - Target")
	fmt.Println(event.ShutdownSchedules)

	err = removeRulesWithTarget(event.Repository, branchName, "stop")
	if err != nil {
		return err
	}

	err = createRulesWithTarget(event.Repository, branchName, event.Branch, "stop", event.ShutdownSchedules)

	return err
}

func removeRulesWithTarget(repository, branch, action string) error {
	fmt.Println("REMOVE SCHEDULES for " + action)

	client := getCloudWatchEventsClient()

	result, err := client.ListRules(nil)
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/removeRuleWithTarget", "operation": "aws/listRules"}, 0)
		return err
	}

	for i := range result.Rules {
		if strings.Contains(*result.Rules[i].Name, "as-"+action+"-"+repository+"-"+branch) {
			fmt.Println("List targets for " + action + " rules ")
			targetResult, err := client.ListTargetsByRule(&cloudwatchevents.ListTargetsByRuleInput{
				Rule: result.Rules[i].Name,
			})
			if err != nil {
				helper.Logger.Log(err, map[string]string{"module": "model/removeRulesWithTarget", "operation": "aws/listTargets"}, 0)
				return err
			}
			helper.Logger.Log(errors.New(fmt.Sprint(targetResult)), map[string]string{"module": "model/removeRulesWithTarget", "operation": "aws/listTargetsResult"}, 4)

			if len(targetResult.Targets) != 0 {
				fmt.Println("Deleting targets")
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
				fmt.Println("Skipping target delete, since there are no targets attached")
			}

			fmt.Println("Deleting rule: " + *result.Rules[i].Name)
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
func createRulesWithTarget(repository, branch, branchRaw, action string, schedule []types.TimeSchedule) error {
	fmt.Println("CREATE NEW SCHEDULES for " + action)

	client := getCloudWatchEventsClient()

	for i := range schedule {
		ruleName := "as-" + action + "-" + repository + "-" + branch + "-" + fmt.Sprint(rand.Intn(9999))

		log.Println("PUT RULE")
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

		log.Println("PUT TARGET")
		target, err := client.PutTargets(&cloudwatchevents.PutTargetsInput{
			Targets: []*cloudwatchevents.Target{
				{
					Arn:   aws.String("arn:aws:lambda:eu-central-1:171842373341:function:auto-staging-scheduler"),
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
