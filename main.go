package main

import (
	"flag"
	"fmt"
	"time"

	"os"

	"github.com/charmbracelet/log"
	"gopkg.in/yaml.v3"
)

var logger *log.Logger

func main() {
	// Setup logger
	logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
		Prefix:          "🎵",
		Level:           log.DebugLevel,
	})

	logger.Info("Welcome to deployer!")

	// Parse flags
	project := flag.String("p", "", "Project to deploy")
	version := flag.String("v", "", "Version to deploy")
	rollback := flag.Bool("r", false, "Rollback to previous version")

	flag.Parse()

	// Validate inputs
	if *project == "" {
		logger.Error("Project is required")
		return
	}

	if *version == "" {
		logger.Error("Version is required")
		return
	}

	// Deploy
	if !*rollback {
		deploy(*project, *version)
	} else {
		rollbackVersion(*project, *version)
	}
}

/**
deployer --project=backend --version=1.2.3
deployer --project=frontend --rollback=true --version=1.2.2

Parse inputs -> get value of project, rollback, version etc
if no rollback, call deploy func

*/

func deploy(project, version string) {
	logger.Debug(fmt.Sprintf("Deploying project: %s version: %s", project, version))
	generateComposeFile(project, version)
}

func rollbackVersion(project, version string) {
}

type ComposeFile struct {
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
	Networks map[string]Network `yaml:"networks"`
	Volumes  map[string]Volume  `yaml:"volumes"`
}

type Service struct {
	Image       string   `yaml:"image"`
	Ports       []string `yaml:"ports"`
	Environment []string `yaml:"environment"`
	Networks    []string `yaml:"networks"`
	Volumes     []string `yaml:"volumes"`
	Deploy      Deploy   `yaml:"deploy"`
}

type Deploy struct {
	Replicas      int           `yaml:"replicas"`
	Resources     Resources     `yaml:"resources"`
	UpdateConfig  UpdateConfig  `yaml:"update_config"`
	RestartPolicy RestartPolicy `yaml:"restart_policy"`
}

type UpdateConfig struct {
	Parallelism int    `yaml:"parallelism"`
	Delay       string `yaml:"delay"`
}

type RestartPolicy struct {
	Condition   string `yaml:"condition"`
	Delay       string `yaml:"delay"`
	MaxAttempts int    `yaml:"max_attempts"`
}

type Resources struct {
	Limits       Limits       `yaml:"limits"`
	Reservations Reservations `yaml:"reservations"`
}

type Limits struct {
	CPUS   string `yaml:"cpus"`
	Memory string `yaml:"memory"`
}

type Reservations struct {
	CPUS   string `yaml:"cpus"`
	Memory string `yaml:"memory"`
}

type Network struct {
	Driver string `yaml:"driver"`
}

type Volume struct {
	Driver string `yaml:"driver"`
}

var ComposeTemplate = `
version : "3.8"

services:
  service1:
    image: image1
    ports: 
      - "5000:5000"
    deploy:
      replicas: 2
      update_config:
        parallelism: 2
        delay: 10s
      restart_policy:
        condition: on-failure
        delay: 3s
      resources:
        limits:
          cpus: "0.5"
          memory: "1G"
        reservations:
          cpus: "0.5"
          memory: "500M"


`

func generateComposeFile(project, version string) {
	// Parse compose file
	var composeFileData ComposeFile
	err := yaml.Unmarshal([]byte(ComposeTemplate), &composeFileData)

	if err != nil {
		logger.Error("Error parsing compose file: ", err)
		return
	}

	logger.Debug("Parsed compose template\n")

	// Update service name
	for serviceName := range composeFileData.Services {
		newServiceName := fmt.Sprintf("%s-%s", project, serviceName)
		composeFileData.Services[newServiceName] = composeFileData.Services[serviceName]
		delete(composeFileData.Services, serviceName)
	}

	// Update image version
	for serviceName, service := range composeFileData.Services {
		service.Image = fmt.Sprintf("%s:%s", service.Image, version)
		composeFileData.Services[serviceName] = service
	}

	// Write new compose file
	newFileContents, err := yaml.Marshal(&composeFileData)
	if err != nil {
		logger.Error("Error writing compose file: ", err)
		return
	}

	logger.Debug("Generated compose file\n")
	fmt.Println(string(newFileContents))

}
