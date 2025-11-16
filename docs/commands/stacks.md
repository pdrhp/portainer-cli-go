# Stacks Command

Manage Docker stacks in Portainer environment.

## Usage

```bash
portainer-go stacks [command]
```

## Available Commands

- `list` - List stacks with optional filters
- `create-swarm-git` - Create a new Swarm stack from a Git repository
- `redeploy` - Redeploy a stack from its Git repository

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

---

## Redeploy Git Command

Redeploy an existing stack by pulling the latest changes from its Git repository.

### Usage

```bash
portainer-go stacks redeploy [stack-id] [flags]
```

### Examples

#### Redeploy with Stack ID as Argument

```bash
portainer-go stacks redeploy 123 --endpoint-id 1 --env DATABASE_URL=prod-db:5432
```

#### Redeploy with Stack ID as Flag

```bash
portainer-go stacks redeploy --stack-id 123 --prune --pull-image
```

#### Redeploy with Git Authentication

```bash
portainer-go stacks redeploy 123 \
  --endpoint-id 1 \
  --repository-username deploy-user \
  --repository-password $GIT_TOKEN \
  --repository-reference-name refs/heads/production
```

#### Redeploy with Environment Variables

```bash
portainer-go stacks redeploy 123 \
  --endpoint-id 1 \
  --env VERSION=v1.2.3 \
  --env ENVIRONMENT=production \
  --env LOG_LEVEL=debug \
  --env API_PORT=8080 \
  --prune
```

#### Redeploy with Complete Configuration

```bash
portainer-go stacks redeploy 123 \
  --endpoint-id 1 \
  --repository-reference-name refs/heads/main \
  --repository-username bot-user \
  --repository-password github-token \
  --env DATABASE_URL=postgres://prod-db:5432/app \
  --env REDIS_URL=redis://cache:6379 \
  --env ENVIRONMENT=production \
  --prune \
  --pull-image \
  --stack-name production-stack
```

#### Interactive Redeploy (Wizard)

```bash
portainer-go stacks redeploy
```

This will start an interactive wizard that guides you through the redeploy configuration.

### Arguments

- `stack-id` - Stack identifier (required, can also be provided via `--stack-id` flag)

### Required Flags

- `--endpoint-id int` - Environment identifier (required)

### Optional Flags

#### Git Configuration

- `--stack-id int` - Stack ID to redeploy (alternative to positional argument)
- `--repository-reference-name string` - Git reference (branch/tag) to pull from
- `--repository-username string` - Username for Git repository authentication
- `--repository-password string` - Password/token for Git repository authentication

#### Stack Configuration

- `--env string` - Environment variables (format: `KEY=value`, can be used multiple times)
- `--prune` - Remove services that are no longer referenced in the compose file
- `--pull-image` - Force pull the latest image even if already present
- `--stack-name string` - Stack name override (Kubernetes only)

### Error Handling

#### Common Errors

##### 400 - Invalid Request

The request parameters are invalid. Check:
- Stack ID format and existence
- Environment variables format (KEY=value)
- Git reference accessibility

##### 403 - Permission Denied

You don't have permission to redeploy this stack. Check your user permissions.

##### 404 - Not Found

Stack with the specified ID doesn't exist.

##### 401 - Authentication Failed

Your authentication token is invalid or expired. Run `portainer-go auth` again.

### Important Notes

- **Environment Variables**: Variables specified with `--env` will override existing environment variables in the stack.
- **Prune Option**: When `--prune` is used, services not present in the compose file will be removed.
- **Pull Image**: Forces a pull of the latest image, even if the tag hasn't changed.
- **Git Authentication**: Only provide credentials if the repository requires authentication for the specified reference.
- **Stack Name**: Only used for Kubernetes stacks, ignored for Docker Swarm stacks.

### CI/CD Integration Examples

#### GitHub Actions

```yaml
name: Redeploy on Push

on:
  push:
    branches: [main]

jobs:
  redeploy:
    runs-on: ubuntu-latest
    steps:
      - name: Redeploy Stack
        run: |
          ./portainer-cli auth --username "${{ secrets.PORTAINER_USER }}" --password "${{ secrets.PORTAINER_PASS }}"
          ./portainer-cli stacks redeploy ${{ secrets.STACK_ID }} \
            --endpoint-id ${{ secrets.ENDPOINT_ID }} \
            --env "COMMIT_SHA=${{ github.sha }}" \
            --env "DEPLOY_TIME=$(date -u +'%Y-%m-%dT%H:%M:%SZ')" \
            --prune \
            --pull-image
```

#### GitLab CI

```yaml
redeploy:production:
  stage: deploy
  script:
    - ./portainer-cli auth --username "$PORTAINER_USER" --password "$PORTAINER_PASS"
    - |
      ./portainer-cli stacks redeploy $STACK_ID \
        --endpoint-id $ENDPOINT_ID \
        --repository-username $GIT_USER \
        --repository-password $GIT_TOKEN \
        --env "CI_COMMIT_SHA=$CI_COMMIT_SHA" \
        --env "CI_PIPELINE_ID=$CI_PIPELINE_ID" \
        --env "DEPLOY_ENV=production" \
        --prune \
        --pull-image
  only:
    - main
```

### Best Practices

1. **Use Environment Variables for Configuration**: Pass runtime configuration via `--env` flags instead of hardcoding in compose files
2. **Enable Prune in Production**: Use `--prune` to ensure removed services are cleaned up
3. **Force Image Pull**: Use `--pull-image` to ensure you get the latest image version
4. **Git Authentication**: Store credentials securely as environment variables
5. **Stack ID Management**: Keep track of stack IDs in your CI/CD variables
6. **Error Handling**: Always check the exit code of redeploy commands in scripts

### Troubleshooting

#### Redeploy Fails with Authentication Error

1. Verify your Portainer authentication token is valid
2. Check if you have permission to redeploy the specified stack
3. Ensure the stack exists and is accessible

#### Git Repository Access Issues

1. Verify repository URL is correct and accessible
2. Check Git credentials if repository requires authentication
3. Ensure the specified branch/tag exists in the repository

#### Stack Not Found

1. Verify the stack ID is correct
2. Ensure you're using the correct endpoint ID
3. Check if the stack was created with Git integration initially
