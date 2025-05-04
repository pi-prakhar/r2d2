package config

type AutoDeployConfig struct {
	Owner       string
	Repo        string
	Tag         string
	Namespace   string
	Deployments []string
	Interval    int
	Token       string
}
