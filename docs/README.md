# Documentation

This directory contains comprehensive documentation for the Portainer CLI.

## Structure

- `commands/` - Detailed documentation for each CLI command
- `api/` - API integration details (future)
- `examples/` - Usage examples and recipes (future)

## Commands

- [auth](commands/auth.md) - Authentication with Portainer
- [config](commands/config.md) - Configuration management
- [stacks](commands/stacks.md) - Stack operations (list and create from Git)

## Contributing to Documentation

When adding new commands or features:

1. Create a new markdown file in `commands/`
2. Follow the existing format with sections for Usage, Examples, Flags, etc.
3. Include error handling and common use cases
4. Add the command to this README

## Generating Documentation

```bash
# Generate command help
go run . [command] --help

# Update this README when adding new commands
```
