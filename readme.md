# Deployer
Deploys a specific tagged image on the machine

# Requirements
## Machine - Linux
## Softwares installed - go, docker, docker compose

# Flow

1. Install the application `go install github.com/MuktadirHassan/deployer`
2. Run `deployer -flag --options`

# Features

1. Versioned deployment 
2. Rollback to a specific version
3. Zero downtime deployment


# Flags
1. -v, --version: Version of the image to deploy
2. -r, --rollback: Rollback to a specific version
3. -h, --help: Help
4. -l, --list: List all the versions available (Not implemented yet)

# Options
1. -c, --config: Path to the configuration file
2. -i, --image: Name of the image to deploy

