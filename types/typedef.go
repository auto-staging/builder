package types

// EnvironmentVariable is the implementation of the TowerAPI EnvironmentVariable schema
type EnvironmentVariable struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Event contains the event body used in the invokation of the Lambda
type Event struct {
	Repository            string                `json:"repository"`
	Branch                string                `json:"branch"`
	Operation             string                `json:"operation"`
	InfrastructureRepoURL string                `json:"infrastructureRepoUrl"`
	CodeBuildRoleARN      string                `json:"codeBuildRoleARN"`
	EnvironmentVariables  []EnvironmentVariable `json:"environmentVariables"`
	Success               int                   `json:"success"`
	ShutdownSchedules     []TimeSchedule        `json:"shutdownSchedules"`
	StartupSchedules      []TimeSchedule        `json:"startupSchedules"`
}

// TimeSchedule is the implementation of the TowerAPI TimeSchedule schema
type TimeSchedule struct {
	Cron string `json:"cron"`
}

// Status struct contains the Environment status
type Status struct {
	Status string `json:"status"`
}

// StatusUpdate struct is used for DynamoDB updates, because the update command requires all json keys to start with ":"
type StatusUpdate struct {
	Status string `json:":status"`
}

// Build struct represents the build step of the CodeBuild buildspec
type Build struct {
	Commands []string
	Finally  []string
}

// Phases struct represents the phases block of the CodeBuild buildspec
type Phases struct {
	Build Build
}

// Buildspec struct represents the general CodeBuild buildspec structure
type Buildspec struct {
	Version string
	Phases  Phases
}
