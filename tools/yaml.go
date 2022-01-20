package tools

import (
	"edgelog/app/model/agent"
	"gopkg.in/yaml.v3"
)

func ReadYaml(configStr string) (agent.Config, error) {
	var config agent.Config
	err := yaml.Unmarshal([]byte(configStr), &config)
	return config, err
}

func WriteYaml(config agent.Config) (string, error) {
	data, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
