package main

import (
	"flag"
	"fmt"
	"os/exec"
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
  service:
    image: programminghero1/prod-neptune-web-backend
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

	logger.Debug("Parsed compose template")

	// Update service names
	for serviceName, service := range composeFileData.Services {
		newServiceName := fmt.Sprintf("%s-%s", project, serviceName)
		logger.Debug(fmt.Sprintf("Updating service name from %s to %s", serviceName, newServiceName))
		delete(composeFileData.Services, serviceName)
		composeFileData.Services[newServiceName] = service
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

	// Write to docker-compose file
	runShellCmd(fmt.Sprintf("echo '%s' > docker-compose.yml", string(newFileContents)))
	runShellCmd(fmt.Sprintf("docker stack deploy -c docker-compose.yml %s", project))

}

func runShellCmd(command string) {
	logger.Debugf("Running command: %s", command)

	// Run command
	cmd := exec.Command("sh", "-c", command)

	out, err := cmd.CombinedOutput()

	if err != nil {
		logger.Errorf(fmt.Sprint(err) + ": \n" + string(out))
		return
	}

	logger.Debugf("Command output: %s", string(out))
}
