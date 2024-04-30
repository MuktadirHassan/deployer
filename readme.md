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

# Must have
.env file with necessary environment variables
deployer.yaml file with necessary configurations


```yaml: deployer.yaml # Path: deployer.yaml
name: Project Name
projects:
    - name: Project 1
      env_path: .env
      configs:
        - name: Config 1
          content: |
            key: value
    - name: Project 2
      env_path: .env


```

# docker swarm setup
```bash
docker swarm init
docker stack deploy -c docker-compose.yml <stack_name>
```

# docker-compose.yml
```yaml: docker-compose.yml
version: '3.7'

services:
  <service_name>:
    image: <image_name>
    ports:
      - "8080:8080"
    environment:
      - key=value
    networks:
      - <network_name>
    deploy:
      replicas: 1
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
        window: 120s
      update_config:
        parallelism: 2
        delay: 10s
      placement:
        constraints:
          - node.role == worker
      resources:
        limits:
          cpus: '0.50'
          memory: 50M
        reservations:
          cpus: '0.25'
          memory: 20M`

networks:
    <network_name>:
        driver: overlay
    
```
