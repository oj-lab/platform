# oj-lab-services

Currently the backend server for OJ Lab.

## Development

🌟 Accept VSCode extension recommandation for complete experience 🌟

For service development, we don't want to make it too complex.
Using VSCode on either Win/*nix System are avaliabe, try using the Makefile/Dockerfile in the repository.

## Serve Frontend

Use `make get-front` to get the frontend dist codes.
You should also set `service.serve_front` to `true` in config file
(see [override.example.toml](config/override.example.toml) for more information)

### WARNING

You should close `remote.autoForwardPorts` if you are using online VSCode
