# Changelog

## [0.2.0] - 2026-05-24
### Fixed
- Added missing `_ "github.com/lib/pq"` import to `delete_user` and `get_user` services (Postgres driver was not registered).
- Fixed `godotenv.Load()` order in `add_user`, `get_users`, `update_user` — now called before `nats.Connect()` so `NATS_URL` is properly loaded.
- `update_user` service was a copy of `add_user` — rewritten with actual `UPDATE` SQL logic.
- Standardized NATS subject names across all services (`users.add`, `users.get`, `users.list`, `users.update`, `users.delete`).
- Removed committed `.env` files from git tracking and added `.gitignore`.
- Created `.env.example` files for all services.
- Cleaned up verbose comments across all source files.
- Updated test files to match new subject names and fixed `update_user` test.
- Updated client files to match new subject names.

## [0.1.0] - 2025-07-28
### Added
- Initial microservices structure: `add_user`, `delete_user`, `get_user`, `get_users`, `update_user`.
- NATS microservice integration with request-reply pattern.
- PostgreSQL persistence for user CRUD operations.
- `run_lint.sh` script for automated linting.
- Unit and integration tests for services.
