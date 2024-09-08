# OJ Lab Platform

![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/oj-lab/platform/build-and-test.yaml?logo=github&label=Tests)
[![codecov](https://codecov.io/gh/oj-lab/platform/graph/badge.svg?token=T5P2KP7VTL)](https://codecov.io/gh/oj-lab/platform)
![Codespace Supported](https://img.shields.io/badge/Codespace_Supported-000000?style=flat&logo=github)

Central service for OJ Lab, supporting distributed deployment.

## Development

> 🌟 Accept VSCode extension recommandation for complete experience.

### Before you start

OJ Lab Platform depends on several foundational services, including:

- PostgreSQL (or other SQL database in the future) for data storage
- Redis for caching & session management
- MinIO (or other S3 like storage) for file storage
- ClickHouse for analytics (currently not developed)
- [Judger](https://github.com/oj-lab/judger) for judging

This project provides a Makefile to help you quickly set up dependencies & other optional choices.
Run `make setup-dependencies` to start these services and load the initial data.

### Launch from VSCode

Launch the programs with VSCode launch configurations is the most recommended way.
It will automatically set the environment and run the program in debug mode.

### Run Judger

There is a `judger` service in the project's `docker-compose.yml`.
It won't start from the `make setup-dependencies` command by default
(since it takes time to let MinIO & PostgreSQL start up).

Run `docker-compose up -d judger` to start the judger service.

### Manage DB data

We provide `adminer` to access PostgreSQL data.
You can optionally run `docker-compose up -d adminer` to start the service,
or just continue to use your own database management tool.

> Remember to set the type of the database to `PostgreSQL` when login to adminer.

## Troubleshooting

go bin not included in PATH

```bash
# You should change .zshrc to .bashrc if you are using bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
```
