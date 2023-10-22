# oj-lab-services

Currently the backend server for OJ Lab.

## Development

🌟 Accept VSCode extension recommandation for complete experience 🌟

For service development, we don't want to make it too complex.
Using VSCode on either Win/*nix System are avaliabe, try using the Makefile/Dockerfile in the repository.

### Run Processes

It's more recommended to use 2 terminals to run the services and background workers seperately.
In this way you will get a better experience of debugging.
(background workers will produce a lot of logs, since they are running in a loop)

```bash
# Terminal 1
make run-server
# Terminal 2
make run-background
```

Alternatively, you can use `make run-all` to run all the processes in one terminal.

## Serve Frontend

Use `make get-front` to get the frontend dist codes.
You should also set `service.serve_front` to `true` in config file
(see [override.example.toml](config/override.example.toml) for more information)

### WARNING

You should close `remote.autoForwardPorts` if you are using online VSCode
