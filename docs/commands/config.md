# Config Command

Manage Portainer CLI global configuration.

## Usage

```bash
portainer-go config [command]
```

## Available Commands

- `set` - Set configuration values
- `get` - Get configuration values

## Examples

### Set Server URL

```bash
portainer-go config set server-url https://portainer-go.company.com
```

### Set Credentials

```bash
portainer-go config set username admin
portainer-go config set password mypassword
```

### View Current Configuration

```bash
portainer-go config get
```

### Get Specific Value

```bash
portainer-go config get server-url
```

## Configuration Keys

- `server-url` - Portainer server URL
- `username` - Default username for authentication
- `password` - Default password for authentication
- `api-key` - API key for authentication (alternative to username/password)

## Configuration File

Configuration is stored in `~/.portainer-go-cli/config.yaml`

Example:
```yaml
server_url: https://portainer-go.company.com
username: admin
password: mypassword
token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Priority Order

1. Command-line flags
2. Configuration file values
3. Interactive prompts

## Security

- Passwords are masked in output
- JWT tokens are stored securely
- Configuration file permissions are restricted
