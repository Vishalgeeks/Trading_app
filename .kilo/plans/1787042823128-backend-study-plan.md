# Backend Codebase Study Plan - Mehndi Booking Application

## Project Overview
**Go 1.26** + **PostgreSQL (pgx/v5)** + **Gorilla Mux** - REST API for henna/mehndi booking system

---

## 1. Architecture & Layer Structure

```
backend/
├── cmd/
│   ├── server/main.go      # App entry point, dependency wiring, HTTP server
│   └── seed/               # Database seeding (not explored)
├── internal/
│   ├── auth/               # Authentication (JWT, bcrypt, register/login)
│   ├── booking/            # Booking domain (CRUD, slots, conflicts, notifications)
│   ├── category/           # Category domain (CRUD, slug generation)
│   ├── design/             # Design domain (CRUD, image upload, search)
│   ├── availability/       # Admin availability slots (day-of-week based)
│   ├── favorite/           # User favorites (simple CRUD)
│   ├── notification/       # Notification system (booking events)
│   ├── user/               # User profile management
│   ├── middleware/         # JWT auth, CORS, role-based access
│   ├── router/             # HTTP route definitions (all endpoints)
│   ├── config/             # Env configuration loading
│   ├── database/           # PostgreSQL connection pool
│   └── handler/            # Shared HTTP helpers (health check)
├── migrations/             # 8 SQL migration files (up/down)
└── go.mod / go.sum         # Dependencies
```

### Layer Responsibilities

| Layer | Files | Responsibility |
|-------|-------|----------------|
| **Entry Point** | `cmd/server/main.go` | Config → DB Pool → Repositories → Services → Handlers → Router → HTTP Server |
| **Config** | `internal/config/config.go` | Load .env, validate required vars (DATABASE_URL, JWT_SECRET) |
| **Database** | `internal/database/postgres.go` | Create pgxpool, ping test |
| **Repository** | `*/repository.go` | Pure SQL queries, scanning rows → structs, error mapping (pgError codes) |
| **Service** | `*/service.go` | Business logic, validation, orchestration, transactions, cross-repo calls |
| **Handler** | `*/handler.go` | HTTP decoding, auth extraction, service calls, JSON responses, status codes |
| **Router** | `internal/router/router.go` | Mux routes, middleware chains (RequireAuth, RequireRole) |
| **Middleware** | `internal/middleware/` | JWT validation, role checks, CORS |
| **Model** | `*/model.go` | Structs, request/response DTOs, conversion methods (ToResponse) |

---

## 2. File-by-File Study Order (Bottom-Up)

### Phase 1: Foundation (Shared Infrastructure)
1. **`internal/config/config.go`** - Config struct, Load(), validation
2. **`internal/database/postgres.go`** - NewPool(), context timeout, ping
3. **`internal/middleware/auth.go`** - JWT parsing, claims extraction, context key, RequireAuth, RequireRole
4. **`internal/middleware/cors.go`** - CORS setup from FRONTEND_URL
5. **`internal/router/router.go`** - All routes mapped, middleware chains
6. **`cmd/server/main.go`** - Full dependency graph wiring

### Phase 2: Auth Domain (Independent)
1. **`internal/auth/model.go`** - Request/Response DTOs, UserResult
2. **`internal/auth/service.go`** - Register, Login, JWT generation, bcrypt, validation
3. **`internal/auth/handler.go`** - HTTP handlers for register, login, me, change-password, logout

### Phase 3: User Domain (Depends on Auth)
1. **`internal/user/model.go`** - User, Role (CLIENT/ADMIN), UpdateProfileRequest
2. **`internal/user/repository.go`** - CRUD + ListUsersByRole (for admin notifications)
3. **`internal/user/service.go`** - Validation, password hashing via auth.HashPassword
4. **`internal/user/handler.go`** - Profile GET/PATCH, UpdateMyProfile (auth required)

### Phase 4: Category Domain (Independent)
1. **`internal/category/model.go`** - Category, Create/Update requests, slug helpers
2. **`internal/category/repository.go`** - CRUD, ListCategories(activeOnly), CountActiveDesignsInCategory
3. **`internal/category/service.go`** - Validation, slug normalization, deactivation check
4. **`internal/category/handler.go`** - Public (List, Get) + Admin (Create, Update, Delete)

### Phase 5: Design Domain (Depends on Category)
1. **`internal/design/model.go`** - Design, Category embedding, Create/Update requests
2. **`internal/design/repository.go`** - CRUD, ListDesigns (pagination, filter, search), slug validation
3. **`internal/design/service.go`** - Category validation (active), slug/price/duration validation
4. **`internal/design/upload.go`** - Multipart file save (random hex name, /uploads/)
5. **`internal/design/handler.go`** - Public + Admin (multipart for image upload)

### Phase 6: Availability Domain (Independent)
1. **`internal/availability/model.go`** - Availability (day_of_week 0-6), requests
2. **`internal/availability/repository.go`** - CRUD, GetAvailabilityForDay (active only)
3. **`internal/availability/service.go`** - Time validation, day range check
4. **`internal/availability/handler.go`** - Admin CRUD

### Phase 7: Booking Domain (Core - Depends on All Above)
1. **`internal/booking/model.go`** - Booking, responses (client/admin), Slot, requests
2. **`internal/booking/repository.go`** - Complex queries:
   - Create/Get/List (with JOINs for design/user names)
   - CheckBookingOverlap (critical for concurrency)
   - Paginated lists with filters
   - Upcoming/History separation
   - Admin stats (conditional aggregation)
   - Scan helpers (scanBooking, scanBookingWithDetails)
3. **`internal/booking/service.go`** - **Most complex**:
   - CreateBooking: availability check → design duration → **serializable TX** → overlap check → insert → notifications
   - Slot generation (30-min intervals, overlap check per slot)
   - Status transitions (PENDING→CONFIRMED→COMPLETED, PENDING/CONFIRMED→CANCELLED)
   - Notifications (client + admins) for create, confirm, complete, cancel
4. **`internal/booking/handler.go`** - Client + Admin endpoints
5. **`internal/booking/slot_handler.go`** - Public GET /api/availability/slots

### Phase 8: Supporting Domains
1. **`internal/favorite/`** - Simple CRUD (user+design unique)
2. **`internal/notification/`** - Types (BOOKING_CREATED, CONFIRMED, etc.), CRUD, unread count, mark read

---

## 3. Key Data Flows

### Booking Creation Flow
```
Client POST /api/bookings
  → Handler.CreateBooking (extract user from JWT)
  → Service.CreateBooking(userID, req)
     → Validate date/time format
     → Check admin availability for day-of-week
     → Get design duration
     → Calculate endTime = startTime + duration
     → BEGIN SERIALIZABLE TRANSACTION
     → CheckBookingOverlap in TX (design, date, start, end)
     → INSERT booking (PENDING)
     → COMMIT
     → Async: sendBookingCreatedNotifications (client + all admins)
  → Return BookingResponse
```

### Slot Calculation Flow
```
GET /api/availability/slots?design_id=X&date=Y
  → SlotHandler.GetAvailableSlots
  → Service.CalculateAvailableSlots
     → Get admin availability for day-of-week
     → Get design duration
     → For each availability window:
        → Generate 30-min slots
        → CheckBookingOverlap for each slot
        → Return available slots
```

### Status Transition Flow
```
Admin PATCH /api/admin/bookings/{id}/status
  → Handler.AdminUpdateBookingStatus
  → Service.UpdateBookingStatus
     → Validate transition (PENDING→CONFIRMED/CANCELLED, CONFIRMED→COMPLETED/CANCELLED)
     → UpdateBookingStatus in repo
     → sendStatusChangeNotification (confirm→notifyConfirmation, complete→notifyCompletion)
```

---

## 4. Test Files Purpose (`*_test.go`)

| Test File | Purpose |
|-----------|---------|
| `booking_test.go` | Integration tests with real PostgreSQL: CreateBooking, Overlap detection, Concurrent booking (race test), Status transitions, Cancel |
| `handler_test.go` (category, user, auth) | HTTP handler tests (require running server) |
| `repository_test.go` (category, user, design) | Repository-level tests with real DB |
| `service_test.go` (category, user, design) | Service logic tests with real DB |
| `password_test.go` | bcrypt hash/verify tests |
| `notification_test.go` | Notification service tests |

**Why separate test files?**
- Go convention: `*_test.go` runs with `go test`
- Separation by layer: repository vs service vs handler
- `TestMain` in booking_test.go sets up/tears down test data
- Tests use real PostgreSQL (not mocks) - requires local DB

---

## 5. Cross-Domain Dependencies (Wiring in main.go)

```
BookingService needs:
  - BookingRepository
  - AvailabilityRepository (for slots)
  - NotificationRepository (for notifications)
  - UserRepository (for admin list in notifications)
  - DesignGetter (interface for design duration)

DesignService needs:
  - DesignRepository
  - CategoryGetter (interface for category validation)

AuthService needs:
  - UserCreator (interface implemented by userRepo adapter in main.go)

NotificationService needs:
  - NotificationRepository
  - UserRepository (for ListUsersByRole)
```

---

## 6. Database Schema (from migrations)

| Table | Key Columns | Constraints |
|-------|-------------|-------------|
| users | id, name, email, phone, password_hash, role, avatar_url, is_active | email UNIQUE, role ENUM(CLIENT,ADMIN) |
| categories | id, name, slug, description, is_active | name/slug UNIQUE |
| designs | id, category_id, name, slug, description, image_url, price, duration_minutes, is_active | slug UNIQUE, category_id FK |
| admin_availability | id, day_of_week (0-6), start_time, end_time, is_active | |
| bookings | id, user_id, design_id, booking_date, start_time, end_time, status, notes | status ENUM, overlap EXCLUDE constraint (migration 006) |
| notifications | id, user_id, type, title, message, booking_id, is_read, read_at | booking_id FK nullable |
| favorites | id, user_id, design_id, created_at | (user_id, design_id) UNIQUE |

---

## 7. Study Approach Recommendations

### For Each Package (Repeat Pattern):
1. **Read `model.go`** - Understand data structures, DTOs, response shapes
2. **Read `repository.go`** - SQL queries, scan functions, error mapping
3. **Read `service.go`** - Business rules, validation, transaction boundaries
4. **Read `handler.go`** - HTTP mapping, auth context, error → status code
5. **Trace in `router.go`** - Which handler method maps to which route + middleware
6. **Check `main.go`** - How dependencies are instantiated and injected

### Critical Deep Dives:
1. **Booking Service** - Transaction handling, overlap prevention, notification orchestration
2. **Auth Middleware** - JWT claims, context propagation, role enforcement
3. **Repository Scan Helpers** - How JOIN results map to nested structs
4. **Slot Generation** - 30-min intervals, overlap checking per slot

### Run Tests to Understand Behavior:
```bash
cd backend
go test ./internal/booking/... -v      # Booking tests (needs PostgreSQL)
go test ./internal/auth/... -v         # Auth tests
go test ./internal/category/... -v     # Category tests
```

---

## 8. Open Questions for Clarification

1. **Migration 006** - "booking_overlap_constraint" - Is this a PostgreSQL EXCLUDE constraint using `btree_gist`? Check the SQL.
2. **Image Upload** - Design handler saves to `./uploads` and serves via `http.FileServer` in main.go. Is this production-ready or dev-only?
3. **Notification Types** - 6 types defined. Are all used? Check service calls.
4. **Concurrency** - Booking uses `SERIALIZABLE` isolation. Any deadlock handling?
5. **Token Revocation** - Logout is no-op. JWTs valid until expiry. Is this intentional?

---

## 9. Next Steps

1. Start with **Phase 1** (Foundation) to understand wiring
2. Pick **one domain** (e.g., Category - simplest) and trace model→repo→service→handler→router
3. Then tackle **Booking** (most complex) with focus on service.go
4. Run tests to see expected behaviors
5. Check migration 006 for overlap constraint implementation

---

*This plan covers 100% backend-only. Frontend integration points noted only where backend serves frontend (CORS, uploads, API contracts).*