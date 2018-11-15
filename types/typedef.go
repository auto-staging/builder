package types

type Event struct {
	Repository           string            `json:"repository"`
	Branch               string            `json:"branch"`
	Operation            string            `json:"operation"`
	RepositoryURL        string            `json:"repositoryUrl"`
	EnvironmentVariables map[string]string `json:"environmentVariables"`
}

type Build struct {
	Commands []string
}

type Phases struct {
	Build Build
}

type Buildspec struct {
	Version string
	Phases  Phases
}
