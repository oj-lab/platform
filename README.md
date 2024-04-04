# OJ Lab Platform

![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/OJ-lab/oj-lab-platform/build-and-test.yaml?logo=github&label=Tests)
![Codespace Supported](https://img.shields.io/badge/Codespace_Supported-000000?style=flat&logo=github)

## Development

> 🌟 Accept VSCode extension recommandation for complete experience.

### Before you start

OJ Lab Platform depends on several foundational services, including:

- PostgreSQL (or other SQL database in the future)
- Redis
- MinIO
- [Judger](https://github.com/OJ-lab/judger)

This project provides a Makefile to help you quickly set up dependencies & other optional choices.
Run `make setup-dependencies` to start these services and load the initial data.

### Launch from VSCode

Launch the programs with VSCode launch configurations is the most recommended way.
It will automatically set the environment and run the program in debug mode.

### Optional

There is also some optional approach you may want to use.

#### Serve frontend

Use `make get-front` to get the frontend dist codes.
> In development config by default, it points to the postion where the frontend codes are located,
> so you will automatically get the frontend view when you start the program.

#### Generate Swagger docs

Use `make gen-swagger` to generate swagger docs.

#### Set environment variables

The following environment variables are available to modify the behavior of the program:

- OJ_LAB_SERVICE_ENV: The environment of the service, default to `development`
- OJ_LAB_WORKDIR: Directly set the path of the workdir, application will automatically locate the workdir if not set
(it will be set to `workdirs/<service_env>` in this project by default)

#### Manage DB data

Along with the `make setup-dependencies`, we provide `adminer` to access PostgreSQL & MinIO with its web interface.

> Remember to set the type of the database to `PostgreSQL` when login to adminer.

## Troubleshooting

go bin not included in PATH

```bash
# You should change .zshrc to .bashrc if you are using bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
```
