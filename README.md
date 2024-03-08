# oj-lab-services

Currently the backend server for OJ Lab.

![Preview](oj-lab-preview.gif)

## Development

🌟 Accept VSCode extension recommandation for complete experience 🌟

For service development, we don't want to make it too complex.
Using VSCode on either Win/*nix System are avaliabe, try using the Makefile/Dockerfile in the repository.

### Serve Frontend

Use `make get-front` to get the frontend dist codes.
You should also set `service.serve_front` to `true` in config file
(see [override.example.toml](environment/configs/override.example.toml) for more information)

### init database

Use `make test` to init data both for **Postgresql**(database) and **minio**(problem-package storage).

while `make setup-db` for **no data** but just table init.

### Run Processes

```bash
make run
```

### WARNING

You should close `remote.autoForwardPorts` if you are using online VSCode

### Troubleshooting

go bin not included in PATH

```bash
# You should change .zshrc to .bashrc if you are using bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
```