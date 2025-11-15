# Stacks Command

Manage Docker stacks in Portainer environment.

## Usage

```bash
portainer-go stacks [command]
```

## Available Commands

- `list` - List stacks with optional filters

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
