# Changelog

## [0.1.0] - 2025-07-28
### Added
- Implemented `run_lint.sh` script to automate `go mod tidy` and `golangci-lint` execution across all microservices.
- Added `.golangci.yml` configuration files for linting in each microservice.
- Enhanced error handling in `delete_user` microservice by checking return values of `req.RespondJSON`.

### Fixed
- Resolved linting issues in `delete_user` microservice.

### Notes
- Initial version of the project with basic microservices structure and linting setup.
