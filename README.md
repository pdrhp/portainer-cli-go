# Portainer Go CLI

A command-line interface for managing Portainer stacks and resources in CI/CD pipelines.

## Features

- 🔐 Authentication with Portainer API
- ⚙️ Global configuration management
- 📋 List stacks with filtering
- 📊 Multiple output formats (table, JSON, YAML)
- 🤖 Interactive wizards for user-friendly experience

## Installation

### From Source

```bash
git clone https://github.com/pdrhp/portainer-go-cli.git
cd portainer-go-cli/app
go build -o portainer-cli .
```

### Using Docker

```bash
docker build -t portainer-cli .
docker run --rm portainer-cli --help
```

## Quick Start

### 1. Configure Portainer server

```bash
./portainer-cli config set server-url https://your-portainer-instance.com
```

### 2. Authenticate

```bash
# Interactive authentication
./portainer-cli auth

# Or with flags
./portainer-cli auth --username admin --password yourpassword
```

### 3. List stacks

```bash
# List all stacks (table format)
./portainer-cli stacks list

# List in JSON format (perfect for CI/CD)
./portainer-cli stacks list --output json

# Filter by endpoint
./portainer-cli stacks list --endpoint-id 1

# Filter by Swarm cluster
./portainer-cli stacks list --swarm-id your-swarm-id
```

## Development

```bash
# Run tests
go test ./...

# Build
go build -o portainer-cli .

# Run
./portainer-cli --help
```

## Available Commands

- `auth` - Authenticate with Portainer server
- `config` - Manage CLI configuration
- `stacks list` - List stacks with optional filters
- `stacks create-swarm-git` - Create a Swarm stack from a Git repository
- `stacks redeploy` - Redeploy a stack from its Git repository

## Examples for CI/CD

```bash
# Get stacks in JSON format for processing
STACKS_JSON=$(./portainer-cli stacks list --output json)

# Check if a specific stack exists
./portainer-cli stacks list --output json | jq '.[] | select(.Name == "my-stack")'

# Count total stacks
./portainer-cli stacks list --output json | jq '. | length'
```

### Create Swarm Stack from Git

```bash
# Interactive creation (wizard)
./portainer-cli stacks create-swarm-git

# Create with flags
./portainer-cli stacks create-swarm-git \
  --name my-stack \
  --repository-url https://github.com/user/repo \
  --swarm-id jpofkc0i9uo9wtx1zesuk649w \
  --endpoint-id 1

# Create with auto-update (GitOps)
./portainer-cli stacks create-swarm-git \
  --name my-stack \
  --repository-url https://github.com/user/repo \
  --swarm-id jpofkc0i9uo9wtx1zesuk649w \
  --endpoint-id 1 \
  --auto-update-interval 1h \
  --auto-update-webhook my-webhook-id
```

### Redeploy Stack from Git

```bash
# Interactive redeploy (wizard)
./portainer-cli stacks redeploy

# Redeploy with flags
./portainer-cli stacks redeploy 123 \
  --endpoint-id 1 \
  --repository-username user \
  --repository-password pass \
  --env DATABASE_URL=prod-db:5432 \
  --prune \
  --pull-image

# Redeploy with new environment variables
./portainer-cli stacks redeploy 123 \
  --endpoint-id 1 \
  --env VERSION=v1.2.3 \
  --env ENVIRONMENT=production
```