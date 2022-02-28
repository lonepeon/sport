# Sport

This application is in charge of storing GPX files for my activities (i.e. running).

All the supported features but also the in progress and upcoming work is listed in the [ROADMAP file](./ROADMAP.md) file.

## Development guidelines

- Everything is done on specific short-lived branches.
  The branches are automatically **rebased to master** if the tests pass.
  It means **all commits should be meaningful and valid**.

- All code should come with its set of tests:
  - Unit test all functionalities, using mocking to only test what is currently at stake
  - Integration test when the feature touches the persistence
  - End-to-end when a new feature is added

## Development tools

- Unit tests: `make test-unit`
- Code quality: `make test-format test-lint test-security`
- Integration tests: `make test-integration` (run tests against `sqlite` or `minio`)
- End-to-end tests: `make test-acceptance`
  Be careful, they will wipe all docker-compose data to always start with a clean state 

## Start the stack locally

- Start `minio` (S3 compatible local replacement) and a `mapbox` mock using `docker-compose up`
- Load the required environment variables.
  They are all listed in the [`Config struct defined in main.go`](./main.go)
- Start the binary `go run .`.
  This will compile and start the application, executing migrations if needed
