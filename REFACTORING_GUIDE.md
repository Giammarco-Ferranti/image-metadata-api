# Step-by-Step Refactoring Guide

This guide will help you refactor your codebase to match the clean architecture described in the architecture document.

## Overview

The refactoring is broken into 10 phases. Each phase builds on the previous one, so follow them in order. After each phase, you should be able to compile and run tests.

---

## Phase 1: Create Configuration Layer

**Goal**: Extract all configuration logic from `cmd/main.go` into a dedicated `pkg/config/` package.

### Steps:

1. **Create `pkg/config/config.go`**

   - Define a `Config` struct with fields:
     - `Port string`
     - `DatabaseURL string` (or separate DB fields: `DBUser`, `DBPassword`, `DBName`, `DBHost`, `DBPort`)
     - `APIKey string`
   - Add struct tags for environment variable names if needed

2. **Create `LoadConfig()` function**

   - Load `.env` file using `godotenv`
   - Read environment variables
   - Validate required fields (return error if missing)
   - Set defaults where appropriate (e.g., `PORT` defaults to "8080")
   - Return `*Config` and `error`

3. **Update `cmd/main.go`**

   - Replace all `os.Getenv()` calls with `config.LoadConfig()`
   - Use the config struct throughout

4. **Test**: Run the application to ensure it still starts correctly

**Checkpoint**: Application should compile and start with the same behavior as before.

---

## Phase 2: Create Domain Layer - Entities

**Goal**: Move domain entities to `pkg/domain/` without database concerns.

### Steps:

1. **Create `pkg/domain/image.go`**

   - Copy `ImageProcess` struct from `pkg/models/models.go`
   - **Remove** all GORM tags (`gorm:"..."`)
   - **Remove** GORM-specific fields if any
   - Keep business fields: `ID`, `CreatedAt`, `UpdatedAt`, `Url`, `Status`, `Width`, `Height`, `Format`
   - Change `sql.NullInt16` to `*int16` (nullable int)
   - Change `sql.NullString` to `*string` (nullable string)
   - Rename to `Image` (domain entities should be simple names)

2. **Add domain behavior methods to `Image`**

   - `Validate()` error - validates URL and business rules
   - `CanDelete()` bool - checks if image can be deleted (status != "in process")
   - `MarkAsProcessing()` - transitions status to "in process"
   - `MarkAsDone(width, height int16, format string)` - transitions to "done" with metadata
   - `MarkAsFailed()` - transitions to "failed"
   - `IsPending()` bool - checks if status is "pending"

3. **Create domain status constants** (optional but recommended)

   - `const StatusPending = "pending"`
   - `const StatusInProcess = "in process"`
   - `const StatusDone = "done"`
   - `const StatusFailed = "failed"`

4. **Keep `pkg/models/models.go` for now** (we'll remove it later)

**Checkpoint**: Domain entity should compile independently. You can write a simple test to verify behavior methods work.

---

## Phase 3: Define Domain Interfaces

**Goal**: Define repository interfaces and Store interface that the domain layer needs.

### Steps:

1. **Create `pkg/domain/store.go`**

   - Define `Store` interface:
     ```go
     type Store interface {
         Atomic(func(Store) error) error
         ImageRepository() ImageRepository
     }
     ```
   - This interface will be implemented by the database layer

2. **Create `pkg/domain/repository.go`**

   - Define `ImageQuerier` interface (read-only):
     ```go
     type ImageQuerier interface {
         FindByID(ctx context.Context, id uuid.UUID) (*Image, error)
         FindAll(ctx context.Context, limit, offset int) ([]*Image, error)
     }
     ```
   - Define `ImageRepository` interface (embeds querier + commands):
     ```go
     type ImageRepository interface {
         ImageQuerier
         Create(ctx context.Context, image *Image) error
         Save(ctx context.Context, image *Image) error
         Delete(ctx context.Context, id uuid.UUID) error
     }
     ```
   - Add `context.Context` to all methods (for cancellation and timeouts)

3. **Create `pkg/domain/commander.go`** (optional for now, but prepare the structure)
   - Define `ImageCommander` interface:
     ```go
     type ImageCommander interface {
         CreateImage(ctx context.Context, url string) (*Image, error)
         // Add other complex operations later
     }
     ```
   - This will encapsulate operations that need transactions

**Checkpoint**: Domain interfaces should compile. No implementations yet.

---

## Phase 4: Create Database Layer

**Goal**: Implement repository interfaces and Store using GORM.

### Steps:

1. **Create `pkg/database/models.go`**

   - Create `ImageModel` struct (GORM model):
     - Copy fields from current `ImageProcess` with GORM tags
     - Keep `sql.NullInt16` and `sql.NullString` for database compatibility
   - Add conversion methods:
     - `ToDomain() *domain.Image` - converts DB model to domain entity
     - `FromDomain(img *domain.Image)` - converts domain entity to DB model

2. **Create `pkg/database/store.go`**

   - Define `store` struct with `*gorm.DB` field
   - Implement `Store` interface:
     - `Atomic(fn func(Store) error) error` - wraps function in GORM transaction
     - `ImageRepository() ImageRepository` - returns repository instance
   - Create `NewStore(db *gorm.DB) Store` factory function

3. **Create `pkg/database/image_repository.go`**

   - Define `imageRepository` struct with `*gorm.DB` field
   - Implement `ImageQuerier` interface:
     - `FindByID` - query DB, convert to domain entity
     - `FindAll` - query DB with pagination, convert to domain entities
   - Implement `ImageRepository` interface:
     - `Create` - insert into DB (convert domain to model first)
     - `Save` - update in DB
     - `Delete` - delete from DB

4. **Create `pkg/database/connection.go`**

   - Move database connection logic from `cmd/main.go`
   - Create `Connect(config *config.Config) (*gorm.DB, error)` function
   - Handle migrations here (move `db.AutoMigrate` call)

5. **Update `cmd/main.go`**
   - Use `database.Connect(config)` instead of direct GORM connection
   - Use `database.NewStore(db)` to create store
   - Remove database connection code

**Checkpoint**: Database layer should compile. You can test repository methods with a simple test.

---

## Phase 5: Refactor API Handlers

**Goal**: Update handlers to use domain layer instead of direct DB access.

### Steps:

1. **Create `pkg/api/requests.go`**

   - Define request DTOs:
     - `CreateImageRequest` struct with `Url string` field

2. **Move response DTOs**

   - Move `ImageProcessResponse` from `pkg/models/models.go` to `pkg/api/responses.go`
   - Rename to `ImageResponse` (or keep as is)
   - Add conversion method: `domain.Image.ToResponse() ImageResponse`

3. **Update `pkg/api/handler.go` (or create if doesn't exist)**

   - Change `Handler` struct:
     ```go
     type Handler struct {
         Querier   domain.ImageQuerier
         Commander domain.ImageCommander  // or Repository for now
     }
     ```
   - Remove `DB *gorm.DB` field

4. **Refactor `handler_add_image.go`**

   - Remove URL validation from handler (move to domain entity `Validate()` method)
   - Remove direct DB access (`h.DB.Create`)
   - Use `h.Commander.CreateImage()` or `h.Repository.Create()`
   - Handle domain errors appropriately

5. **Refactor `handler_get_image.go`**

   - Replace `h.DB.Where(...)` with `h.Querier.FindByID()`
   - Handle `domain.ErrNotFound` type errors

6. **Refactor `handler_get_images.go`**

   - Replace `h.DB.Limit(...).Offset(...)` with `h.Querier.FindAll()`
   - Keep pagination logic in handler (or move to domain if you prefer)

7. **Refactor `handler_delete_image.go`**

   - Use `h.Querier.FindByID()` to get image
   - Use domain entity's `CanDelete()` method instead of checking status in handler
   - Use `h.Repository.Delete()` instead of direct DB access

8. **Update `pkg/api/routes.go`**

   - Update handler initialization to pass querier/commander instead of DB

9. **Update `cmd/main.go`**
   - Create store: `store := database.NewStore(db)`
   - Create querier: `querier := store.ImageRepository()` (or create separate querier)
   - Create commander: `commander := ...` (or use repository for now)
   - Pass to handler: `apiHandler := api.Handler{Querier: querier, Commander: commander}`

**Checkpoint**: All handlers should compile. Test API endpoints to ensure they still work.

---

## Phase 6: Create Authorization Layer

**Goal**: Separate authentication and authorization concerns.

### Steps:

1. **Create `pkg/authz/` directory**

2. **Create `pkg/authz/middleware.go`**

   - Move authentication logic from `pkg/middleware/auth.go`
   - Extract token and add to context (e.g., `context.WithValue(ctx, "apiKey", token)`)
   - Keep it simple for now (just API key validation)

3. **Create `pkg/authz/rules.go`** (optional for now)

   - Prepare structure for future authorization rules
   - For now, you can keep simple API key auth

4. **Update `pkg/api/routes.go`**

   - Change import from `pkg/middleware` to `pkg/authz`
   - Update middleware reference

5. **Remove `pkg/middleware/` directory** (or keep for other middleware like CORS)

**Checkpoint**: Authentication should still work. API endpoints require valid API key.

---

## Phase 7: Implement Commanders

**Goal**: Create command pattern for complex operations that need transactions.

### Steps:

1. **Implement `ImageCommander` in `pkg/database/image_commander.go`**

   - Define `imageCommander` struct with `Store` field
   - Implement `CreateImage(ctx context.Context, url string) (*domain.Image, error)`:
     - Use `store.Atomic()` to wrap in transaction
     - Create domain entity
     - Validate entity
     - Save via repository
     - Return created entity

2. **Update `pkg/api/handler.go`**

   - Ensure `Commander` field uses `ImageCommander` interface

3. **Update handlers to use commanders**

   - `handler_add_image.go` should use `h.Commander.CreateImage()`

4. **Add more commands as needed**:
   - `UpdateImageStatus` - for worker to update status
   - `DeleteImage` - if deletion needs complex logic

**Checkpoint**: Commanders should handle transactions properly. Test that operations are atomic.

---

## Phase 8: Refactor Worker

**Goal**: Update worker to use domain layer instead of direct DB access.

### Steps:

1. **Update `pkg/worker/process_images.go`**

   - Change `StartExtract` signature:
     - Remove `*gorm.DB` parameter
     - Add `domain.ImageQuerier` and `domain.ImageCommander` (or `domain.ImageRepository`) parameters
   - Replace `DB.Where(...)` with `querier.FindAll()` or create a method like `FindPendingImages()`
   - Replace `DB.Model(...).Updates(...)` with domain entity methods and commander

2. **Update `ProcessImage` function**

   - Remove `*gorm.DB` parameter
   - Add repository/commander parameters
   - Use domain entity methods (`MarkAsProcessing()`, `MarkAsDone()`, `MarkAsFailed()`)
   - Use commander or repository to save changes

3. **Update `cmd/main.go`**
   - Pass querier and commander to `worker.StartExtract()` instead of `*gorm.DB`

**Checkpoint**: Worker should compile and process images using domain layer.

---

## Phase 9: Reorganize Directory Structure

**Goal**: Move files to match the target architecture structure.

### Steps:

1. **Move entry point**

   - Create `cmd/server/` directory
   - Move `cmd/main.go` → `cmd/server/main.go`
   - Update any build scripts or Dockerfiles that reference `cmd/main.go`

2. **Clean up old files**

   - Remove `pkg/models/models.go` (domain entities are now in `pkg/domain/`)
   - Remove `pkg/middleware/auth.go` (moved to `pkg/authz/`)

3. **Move test files** (optional, Go convention is to keep tests with code)

   - If you want to follow the architecture doc exactly:
     - Move HTTP test files from `tests/api/` → `test/rest/`
   - Or keep tests next to code (more idiomatic Go)

4. **Update imports**
   - Search and replace all imports that reference old paths
   - Update `go.mod` if module path changed

**Checkpoint**: All imports should resolve. Application should compile and run.

---

## Phase 10: Update Tests

**Goal**: Refactor tests to work with new architecture.

### Steps:

1. **Create mock interfaces** (if using mocks)

   - Create `pkg/domain/mocks.go` or use a mocking library
   - Define mock implementations of `ImageQuerier`, `ImageRepository`, `ImageCommander`

2. **Update handler tests** (`tests/api/*.go`)

   - Replace GORM mocks with repository interface mocks
   - Update handler initialization to use mocks
   - Test domain logic separately from API logic

3. **Update worker tests** (`tests/worker/*.go`)

   - Replace `*gorm.DB` mocks with repository mocks
   - Test worker logic with domain entities

4. **Create domain tests**

   - Test entity behavior methods (`Validate()`, `CanDelete()`, etc.)
   - Test status transitions

5. **Create repository tests**
   - Test database layer implementations
   - Use test database or in-memory database

**Checkpoint**: All tests should pass. Coverage should be maintained or improved.

---

## Additional Considerations

### Error Handling

- Define domain errors in `pkg/domain/errors.go`:
  - `ErrImageNotFound`
  - `ErrInvalidURL`
  - `ErrInvalidStatus`
- Map domain errors to HTTP status codes in API layer

### Migration Strategy

- You can work incrementally: keep old code working while adding new layers
- Test after each phase
- Use feature flags if needed to gradually roll out changes

### Testing Strategy

- Start with unit tests for domain layer (no dependencies)
- Then test repository layer with test database
- Finally test API layer with mocked repositories

---

## Quick Reference: File Structure After Refactoring

```
/
├── cmd/
│   └── server/
│       └── main.go
├── pkg/
│   ├── api/
│   │   ├── handler_add_image.go
│   │   ├── handler_delete_image.go
│   │   ├── handler_get_image.go
│   │   ├── handler_get_images.go
│   │   ├── handler.go
│   │   ├── requests.go
│   │   ├── responses.go
│   │   └── routes.go
│   ├── authz/
│   │   └── middleware.go
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   ├── connection.go
│   │   ├── image_repository.go
│   │   ├── image_commander.go
│   │   ├── models.go
│   │   └── store.go
│   ├── domain/
│   │   ├── image.go
│   │   ├── repository.go
│   │   ├── commander.go
│   │   ├── store.go
│   │   └── errors.go
│   ├── health/
│   │   └── handler_health.go
│   ├── responses/
│   │   └── responses.go
│   └── worker/
│       └── process_images.go
└── test/
    └── rest/
        └── (HTTP test files)
```

---

## Tips for Success

1. **Work incrementally**: Complete one phase before moving to the next
2. **Test frequently**: Run tests after each phase
3. **Use git**: Commit after each successful phase
4. **Keep it simple**: Don't over-engineer. Start with the basics and add complexity as needed
5. **Refer to architecture doc**: Keep the architecture document open for reference

Good luck with your refactoring!
