# Auth Command

Authenticate with Portainer server to obtain a JWT token for subsequent operations.

## Usage

```bash
portainer-go auth [flags]
```

## Examples

### Interactive Authentication

```bash
portainer-go auth
```

This will prompt you to enter your username and password interactively.

### Authentication with Flags

```bash
portainer-go auth --username admin --password mypassword
```

### Authentication with Config Values

If you have previously configured credentials:

```bash
portainer-go config set username admin
portainer-go config set password mypassword
portainer-go auth
```

## Flags

- `--username string` - Username for authentication
- `--password string` - Password for authentication

## Global Flags

- `--server-url string` - Portainer server URL (can be set via config)

## Behavior

1. Validates server URL is configured or provided
2. Attempts authentication with Portainer API
3. Stores JWT token securely in local config file
4. Token is automatically used for subsequent API calls

## Error Handling

- **400/422**: Invalid credentials
- **401**: Authentication failed
- **500**: Server error
- **Connection errors**: Network issues

## Configuration

The authentication token is stored in `~/.portainer-go-cli/config.yaml` and reused automatically until expiration.
