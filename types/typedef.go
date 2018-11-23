package types

type EnvironmentVariable struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Event struct {
	Repository            string                `json:"repository"`
	Branch                string                `json:"branch"`
	Operation             string                `json:"operation"`
	InfrastructureRepoURL string                `json:"infrastructureRepoUrl"`
	CodeBuildRoleARN      string                `json:"codeBuildRoleARN"`
	EnvironmentVariables  []EnvironmentVariable `json:"environmentVariables"`
	Success               int                   `json:"success"`
}

type StatusUpdate struct {
	Status string `json:":status"`
}

type Build struct {
	Commands []string
	Finally  []string
}

type Phases struct {
	Build Build
}

type Buildspec struct {
	Version string
	Phases  Phases
}
