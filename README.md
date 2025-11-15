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
git clone https://github.com/yourorg/portainer-go.git
cd portainer-go
cd app && go build -o portainer-go .
```

### Using Docker

```bash
docker build -t portainer-go app/
docker run --rm portainer-go --help
```

## Quick Start

### 1. Configure Portainer server

```bash
./portainer-go config set server-url https://your-portainer-instance.com
```

### 2. Authenticate

```bash
# Interactive authentication
./portainer-go auth

# Or with flags
./portainer-go auth --username admin --password yourpassword
```

### 3. List stacks

```bash
# List all stacks (table format)
./portainer-go stacks list

# List in JSON format (perfect for CI/CD)
./portainer-go stacks list --output json

# Filter by endpoint
./portainer-go stacks list --endpoint-id 1

# Filter by Swarm cluster
./portainer-go stacks list --swarm-id your-swarm-id
```

## Development

```bash
# Run tests
go test ./...

# Build
cd app && go build -o portainer-go .

# Run
./portainer-go --help
```

## Available Commands

- `auth` - Authenticate with Portainer server
- `config` - Manage CLI configuration
- `stacks list` - List stacks with optional filters

## Examples for CI/CD

```bash
# Get stacks in JSON format for processing
STACKS_JSON=$(./portainer-go stacks list --output json)

# Check if a specific stack exists
./portainer-go stacks list --output json | jq '.[] | select(.Name == "my-stack")'

# Count total stacks
./portainer-go stacks list --output json | jq '. | length'
```

## License

[Your License Here]
