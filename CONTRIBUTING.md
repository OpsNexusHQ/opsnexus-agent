# Contributing to OpsNexus Agent

Use a focused branch and pull request. Explain collector or transport impact, supported Linux behavior, configuration changes, tests, and rollback considerations.

Before opening a pull request:

- run gofmt on changed Go files and go test ./...;
- update README/configuration documentation when environment or CLI behavior changes;
- never commit credentials, host data, binaries, or local configuration;
- coordinate wire-format changes with opsnexus-common, opsnexus-api, and opsnexus-backend.

Maintainers review resource usage, portability, error handling, privacy, and compatibility with the supported backend contract.
