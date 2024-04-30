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

	logger.Debug("Welcome to deployer!")

	// Parse flags
	project := flag.String("p", "", "Project to deploy")
	version := flag.String("v", "", "Version to deploy")
	rollback := flag.Bool("r", false, "Rollback to previous version")

	flag.Parse()

	logger.Debug("________Flags: %v________")
	logger.Debugf("Project: %s", *project)
	logger.Debugf("Version: %s", *version)
	logger.Debugf("Rollback: %t", *rollback)

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
}

type Network struct {
	Driver string `yaml:"driver"`
}

type Volume struct {
	Driver string `yaml:"driver"`
}

func generateComposeFile(project, version string) {

	// Generate compose file with project and version
	composeFile := ComposeFile{
		Version: "3",
		Services: map[string]Service{
			project: {
				Image:       project + ":" + version,
				Ports:       []string{"8080:8080"},
				Environment: []string{"ENV=PROD"},
				Networks:    []string{"default"},
				Volumes:     []string{"./" + project + "/data:/data"},
			},
		},
		Networks: map[string]Network{
			"default": {
				Driver: "bridge",
			},
		},
		Volumes: map[string]Volume{
			"./" + project + "/data": {
				Driver: "local",
			},
		},
	}

	d, err := yaml.Marshal(&composeFile)
	if err != nil {
		logger.Error("Error marshalling compose file: ", err)
		return
	}

	logger.Debug("Compose file: ", string(d))

}
