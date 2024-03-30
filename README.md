# oj-lab-services

![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/OJ-lab/oj-lab-platform/build-and-test.yaml?logo=github&label=Tests)
![Codespace Supported](https://img.shields.io/badge/Codespace_Supported-000000?style=flat&logo=github)

Currently the backend server for OJ Lab.

## Development

🌟 Accept VSCode extension recommandation for complete experience 🌟

For service development, we don't want to make it too complex.
Using VSCode on either Win/*nix System are avaliabe, try using the Makefile/Dockerfile in the repository.

### Serve Frontend

Use `make get-front` to get the frontend dist codes.
You should also set `service.serve_front` to `true` in config file
(see [override.example.toml](environment/configs/override.example.toml) for more information)

### Init Database

Check config in `environment/configs`.

```bash
# To initialize the minio config
cp environment/rclone-minio.example.conf environment/rclone-minio.conf
```

Use `make pkg` to init database with **package** data.

While `make setup-db` for **no data** but just table and user init.

### Run Processes

```bash
make run
```

`make all` do previous things except **config** in one line.

### Troubleshooting

go bin not included in PATH

```bash
# You should change .zshrc to .bashrc if you are using bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
```
