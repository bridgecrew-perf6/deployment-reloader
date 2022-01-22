package model

type DeploymentName string

type  EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type ConfigMapData map[DeploymentName][]EnvVar