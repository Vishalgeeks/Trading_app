# Phase 2 — User & Database Foundation Plan

## Current State
Phase 1 is complete. Codebase has `cmd/server`, `internal/config`, `internal/database`, `internal/handler`, `internal/router`, `internal/middleware`. PostgreSQL via pgxpool. No sqlc, no user module, no migrations beyond `.gitkeep`.

## Objectives
Implement Users module: migration, sqlc queries, domain model, repository, service, handler, routes, tests. No auth/JWT. Only CLIENT and ADMIN roles. No staff/artist concepts.

## Steps

### 1. Add sqlc
- Create `sqlc.yaml` with:
  - `version: "2"`
  - `sql:` pointing to `sql/queries/` and `sql/generated/`
  - `engine: postgresql`
  - `queries` root
  - `schema: migrations/`
- Run `sqlc generate` after migration and queries exist.

### 2. Database Migration
- Create `migrations/001_create_users.up.sql`:
  - `CREATE TYPE user_role AS ENUM ('CLIENT', 'ADMIN');`
  - `CREATE TABLE users (...)` with all specified columns and constraints.
- Create `migrations/001_create_users.down.sql`:
  - `DROP TABLE users;`
  - `DROP TYPE user_role;`
- Verify `migrate up`, `migrate down 1`, `migrate up`.

### 3. sqlc Queries
- Create `sql/queries/users.sql` with:
  - `CreateUser` (`:one`)
  - `GetUserByID` (`:one`)
  - `GetUserByEmail` (`:one`)
  - `UpdateUserProfile` (`:one`)
  - `DeactivateUser` (`:exec`)
  - `ListUsers` (`:many`)
  - `ListUsersByRole` (`:many`)
- Use PostgreSQL enum literal `user_role` in queries.

### 4. Domain Model
- Create `internal/user/model.go`:
  - `type User struct` mirroring table columns.
  - `type Role string` with `const (RoleClient Role = "CLIENT" RoleAdmin Role = "ADMIN")`.
  - Never include `PasswordHash` in API-facing DTOs.

### 5. Repository
- Create `internal/user/repository.go`:
  - Depends on `*pgxpool.Pool` (reuse Phase 1 connection).
  - Implement all 7 methods using sqlc generated code.
  - Map sqlc types to domain `User`.
  - Translate `pgx.ErrNoRows` to sentinel errors (`ErrUserNotFound`).
  - Translate unique-violation `pgconn.PgError.Code == "23505"` to `ErrDuplicateEmail`.

### 6. Service
- Create `internal/user/service.go`:
  - Depends on repository interface.
  - `CreateUser`: validate name/email/role, hash password with bcrypt, call repo.
  - `GetUserByID`, `GetUserByEmail`: pass-through with `ErrUserNotFound`.
  - `UpdateUserProfile`: restrict updates to name, phone, avatar_url; return `ErrUserNotFound`.
  - `DeactivateUser`: set `is_active = false`.
  - `ListUsers`, `ListUsersByRole`: pagination omitted in Phase 2 unless required by plan.
  - Validate email format.
  - Do not expose raw DB errors.

### 7. Handler & Router Integration
- Create `internal/user/handler.go`:
  - `GetUserProfile(w, r, id)` — call service, JSON encode User (no PasswordHash).
  - `UpdateUserProfile(w, r, id)` — decode JSON, call service, return updated User.
- Update `internal/router/router.go`:
  - Add `GET /api/users/{id}`, `PATCH /api/users/{id}`.
  - No auth middleware yet.

### 8. Tests
- **Repository**: integration tests against local DB. Use `t.Parallel`, rollback via transaction, create test user.
- **Service**: unit tests with mock repository (use simple hand-written mock, no heavy framework).
- **Handler**: table-driven HTTP tests using `httptest`.
- **Migration**: simple shell/Go checks for up/down/re-up.
- **Business rules covered**:
  - duplicate email returns conflict error
  - empty name rejected
  - invalid email rejected
  - invalid role rejected
  - profile update blocks role/is_active/password changes
  - GET/PATCH invalid UUID → 400
  - missing user → 404

### 9. Validation
- `go mod tidy` then `go fmt ./...`, `go vet ./...`, `go test ./...`, `go build ./...`
- Verify `GET /health`.
- Verify `GET /api/users/:id` (need seed user).
- Verify `PATCH /api/users/:id`.
- Verify `password_hash` never in response.

## Decisions
- **sqlc package**: `sqlc` in `sql/generated/`.
- **Connection**: reuse existing `*pgxpool.Pool`.
- **Password**: bcrypt via `golang.org/x/crypto/bcrypt` (check if already in go.mod; likely need to add).
- **Testing**: repository = integration against real DB. service/handler = unit with mocks.
- **No pagination** in Phase 2.

## Risks
- PostgreSQL enum usage with sqlc requires correct type casting in queries.
- `DATABASE_URL` must point to existing DB; migration auto-runs not implemented, so `migrate up` must be run manually during setup.
- `golang.org/x/crypto` not currently in go.mod.

## Files to Create
- `sqlc.yaml`
- `sql/queries/users.sql`
- `sql/generated/` (via `sqlc generate`)
- `internal/user/model.go`
- `internal/user/repository.go`
- `internal/user/service.go`
- `internal/user/handler.go`
- `migrations/001_create_users.up.sql`
- `migrations/001_create_users.down.sql`
- `internal/user/*_test.go`
- `migrations_test.go` or equivalent

## Files to Modify
- `go.mod` (add sqlc, bcrypt)
- `internal/router/router.go` (add user routes)
