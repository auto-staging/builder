package types

type Event struct {
	Repository           string            `json:"repository"`
	Branch               string            `json:"branch"`
	Operation            string            `json:"operation"`
	RepositoryURL        string            `json:"repositoryUrl"`
	EnvironmentVariables map[string]string `json:"environmentVariables"`
	Success              int               `json:"success"`
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
