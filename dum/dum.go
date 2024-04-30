package dum

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Name    string
	Content string
}

type ProjectConfig struct {
	Name    string   `yaml:"name"`
	EnvPath string   `yaml:"env_path"`
	Configs []Config `yaml:"configs"`
}

type Project struct {
	Name     string          `yaml:"name"`
	Projects []ProjectConfig `yaml:"projects"`
}

func Foot() {
	project := Project{
		Name: "Project Name",
		Projects: []ProjectConfig{
			{
				Name:    "Project 1",
				EnvPath: ".env",
				Configs: []Config{
					{
						Name:    "Config 1",
						Content: "key: value",
					},
				},
			},
			{
				Name:    "Project 2",
				EnvPath: ".env",
			},
		},
	}

	// Marshal the struct into YAML
	yamlData, err := yaml.Marshal(project)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(string(yamlData))
}
