# Docker Swarm configs

`configs` manages Docker Swarm configs. The existing `config` command manages CLI settings.

Authenticate first using `portainer-go auth`. The command uses the saved JWT token and server URL; `--server-url` overrides the URL.

## Create from a file

```bash
portainer-go configs create nginx-conf-v1 --endpoint-id 1 --file infra/nginx.conf
```

The default output is the created config ID. Use `--output json` or `--output yaml` for structured output containing `ID`.

## Create from stdin

```bash
cat infra/nginx.conf | portainer-go configs create nginx-conf-v1 --endpoint-id 1 --file - --output json
```

Files are sent byte-for-byte, Base64 encoded in the Docker API request. The file must be readable, and the environment ID must be positive. The target environment must provide access to a Swarm manager, and the authenticated user must have permission to create configs. API errors are returned as command failures.

## Reference from a stack

```yaml
services:
  nginx:
    image: nginx:alpine
    configs:
      - source: nginx_conf
        target: /etc/nginx/nginx.conf
configs:
  nginx_conf:
    external: true
    name: nginx-conf-v1
```

Create the config before deploying the stack. Config content is immutable: use a new versioned name when changing the file, then update the stack reference. Creating the same name again returns the Docker API error; this command does not replace or delete existing configs. Keep previous versions if needed for rollback. Use Docker secrets for sensitive data.

## API

`POST /api/endpoints/{endpoint-id}/docker/configs/create`

References: [Portainer Docker proxy](https://docs.portainer.io/api/examples), [Docker API](https://docs.docker.com/reference/api/engine/version/v1.51/), [Swarm configs](https://docs.docker.com/engine/swarm/configs/).

## Validation

The automated suite exercises the public client and CLI against local HTTP test servers. It checks request paths, authentication, encoded content, file/stdin input, invalid arguments and API failures. It does not require or modify a real Portainer/Swarm environment.
