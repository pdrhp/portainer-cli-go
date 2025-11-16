# Stacks Command

Manage Docker stacks in Portainer environment.

## Usage

```bash
portainer-go stacks [command]
```

## Available Commands

- `list` - List stacks with optional filters
- `create-swarm-git` - Create a new Swarm stack from a Git repository

## Examples

### List All Stacks

```bash
portainer-go stacks list
```

Output:
```
ID  NAME        TYPE    STATUS    ENDPOINT  SWARM ID
--  ----        ----    ------    --------  --------
1   web-app     compose running  1         -
2   api-service swarm   running  1         jpofkc0i9uo9wtx1zesuk649w
```

### List in JSON Format

```bash
portainer-go stacks list --output json
```

### Filter by Endpoint

```bash
portainer-go stacks list --endpoint-id 1
```

### Filter by Swarm Cluster

```bash
portainer-go stacks list --swarm-id jpofkc0i9uo9wtx1zesuk649w
```

## Flags

### List Command Flags

- `--endpoint-id int` - Filter stacks by endpoint ID
- `--swarm-id string` - Filter stacks by Swarm cluster ID

### Global Flags

- `--output string` - Output format: table, json, yaml (default "table")
- `--server-url string` - Portainer server URL

## Output Formats

### Table (Default)

Human-readable tabular format with truncated values for better display.

### JSON

Complete JSON representation of all stack data, perfect for CI/CD pipelines.

```json
[
  {
    "Id": 1,
    "Name": "web-app",
    "Type": 1,
    "EndpointId": 1,
    "SwarmId": "",
    "Status": 1,
    "ProjectPath": "/data/compose/1",
    "EntryPoint": "docker-compose.yml"
  }
]
```

### YAML

YAML representation of stack data.

## Status Values

- `1` - Running
- `2` - Stopped
- `3` - Failed

## Type Values

- `1` - Docker Compose
- `2` - Docker Swarm

## Permissions

Only stacks accessible to the authenticated user are shown. Administrators see all stacks, regular users see only stacks they have access to.

---

## Create Swarm Git Command

Create a new Docker Swarm stack by pulling the compose file from a Git repository.

### Usage

```bash
portainer-go stacks create-swarm-git [flags]
```

### Examples

#### Interactive Mode (Wizard)

```bash
portainer-go stacks create-swarm-git
```

This will start an interactive wizard that guides you through all configuration options.

#### Basic Creation with Required Flags

```bash
portainer-go stacks create-swarm-git \
  --name my-stack \
  --repository-url https://github.com/user/repo \
  --swarm-id jpofkc0i9uo9wtx1zesuk649w \
  --endpoint-id 1
```

#### With Custom Compose File and Branch

```bash
portainer-go stacks create-swarm-git \
  --name my-stack \
  --repository-url https://github.com/user/repo \
  --swarm-id jpofkc0i9uo9wtx1zesuk649w \
  --endpoint-id 1 \
  --compose-file docker-compose.prod.yml \
  --repository-reference-name refs/heads/production
```

#### With Environment Variables

```bash
portainer-go stacks create-swarm-git \
  --name my-stack \
  --repository-url https://github.com/user/repo \
  --swarm-id jpofkc0i9uo9wtx1zesuk649w \
  --endpoint-id 1 \
  --env DATABASE_URL=postgres://db:5432/mydb \
  --env API_KEY=secret123
```

#### With Git Authentication

```bash
portainer-go stacks create-swarm-git \
  --name my-stack \
  --repository-url https://github.com/user/private-repo \
  --swarm-id jpofkc0i9uo9wtx1zesuk649w \
  --endpoint-id 1 \
  --repository-username myuser \
  --repository-password mytoken
```

#### With GitOps Auto-Update

```bash
portainer-go stacks create-swarm-git \
  --name my-stack \
  --repository-url https://github.com/user/repo \
  --swarm-id jpofkc0i9uo9wtx1zesuk649w \
  --endpoint-id 1 \
  --auto-update-interval 1h \
  --auto-update-webhook my-webhook-id
```

### Required Flags

- `--name string` - Name of the stack
- `--repository-url string` - URL of the Git repository
- `--swarm-id string` - Swarm cluster identifier
- `--endpoint-id int` - Identifier of the environment

### Optional Flags

#### Git Configuration

- `--compose-file string` - Path to the compose file in the repository (default: `docker-compose.yml`)
- `--repository-reference-name string` - Git reference (branch/tag) (default: `refs/heads/master`)
- `--tlsskip-verify` - Skip TLS verification for Git repository
- `--repository-username string` - Username for Git repository authentication
- `--repository-password string` - Password/token for Git repository authentication

#### Stack Configuration

- `--env string` - Environment variables (format: `KEY=value`, can be used multiple times)
- `--additional-files string` - Additional compose files (can be used multiple times)

#### GitOps Auto-Update

- `--auto-update-interval string` - Auto-update interval (e.g., `1h`, `30m`)
- `--auto-update-webhook string` - Webhook ID for auto-update triggers
- `--auto-update-force-pull-image` - Force pull latest image on auto-update
- `--auto-update-force-update` - Force update even without repository changes

### Error Handling

#### Common Errors

##### 400 - Invalid Request

The request parameters are invalid. Check:
- Swarm ID format
- Repository URL accessibility
- Compose file path exists in repository

##### 409 - Conflict

Stack name or webhook ID already exists. Use a different name or webhook ID.

##### 401/403 - Authentication Failed

Your authentication token is invalid or expired. Run `portainer-go auth` again.

### GitOps Auto-Update

When auto-update is enabled, Portainer will automatically check the Git repository for changes and update the stack accordingly. This enables true GitOps workflows where your Git repository is the source of truth.

#### Auto-Update Interval

The interval defines how often Portainer checks for updates. Examples:
- `30s` - Every 30 seconds
- `5m` - Every 5 minutes
- `1h` - Every hour

#### Webhook Triggers

You can also trigger updates via webhook. The webhook ID must be unique across all stacks.
