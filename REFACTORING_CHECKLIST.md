# Refactoring Checklist

Use this checklist to track your progress through each phase.

## Phase 1: Configuration Layer ✅/❌

- [ ] Created `pkg/config/config.go` with `Config` struct
- [ ] Implemented `LoadConfig()` function
- [ ] Updated `cmd/main.go` to use config package
- [ ] Application compiles and runs
- [ ] **Checkpoint**: Application starts correctly

---

## Phase 2: Domain Entities ✅/❌

- [ ] Created `pkg/domain/image.go` with `Image` entity
- [ ] Removed all GORM tags from entity
- [ ] Changed `sql.NullInt16` to `*int16`
- [ ] Changed `sql.NullString` to `*string`
- [ ] Added `Validate()` method
- [ ] Added `CanDelete()` method
- [ ] Added `MarkAsProcessing()` method
- [ ] Added `MarkAsDone()` method
- [ ] Added `MarkAsFailed()` method
- [ ] Added status constants
- [ ] **Checkpoint**: Domain entity compiles independently

---

## Phase 3: Domain Interfaces ✅/❌

- [ ] Created `pkg/domain/store.go` with `Store` interface
- [ ] Created `pkg/domain/repository.go` with `ImageQuerier` interface
- [ ] Created `pkg/domain/repository.go` with `ImageRepository` interface
- [ ] Created `pkg/domain/commander.go` with `ImageCommander` interface
- [ ] All interfaces use `context.Context`
- [ ] **Checkpoint**: Domain interfaces compile

---

## Phase 4: Database Layer ✅/❌

- [ ] Created `pkg/database/models.go` with `ImageModel` struct
- [ ] Added `ToDomain()` conversion method
- [ ] Added `FromDomain()` conversion method
- [ ] Created `pkg/database/store.go` with Store implementation
- [ ] Implemented `Atomic()` method with transactions
- [ ] Created `pkg/database/image_repository.go`
- [ ] Implemented `ImageQuerier` interface
- [ ] Implemented `ImageRepository` interface
- [ ] Created `pkg/database/connection.go`
- [ ] Moved database connection logic from `cmd/main.go`
- [ ] Updated `cmd/main.go` to use database package
- [ ] **Checkpoint**: Database layer compiles and works

---

## Phase 5: Refactor API Handlers ✅/❌

- [ ] Created `pkg/api/requests.go` with request DTOs
- [ ] Moved `ImageProcessResponse` to `pkg/api/responses.go`
- [ ] Added conversion method `domain.Image.ToResponse()`
- [ ] Updated `Handler` struct to use querier/commander
- [ ] Removed `DB *gorm.DB` from Handler
- [ ] Refactored `handler_add_image.go`
- [ ] Refactored `handler_get_image.go`
- [ ] Refactored `handler_get_images.go`
- [ ] Refactored `handler_delete_image.go`
- [ ] Updated `pkg/api/routes.go`
- [ ] Updated `cmd/main.go` to wire dependencies
- [ ] **Checkpoint**: All handlers compile and API works

---

## Phase 6: Authorization Layer ✅/❌

- [ ] Created `pkg/authz/` directory
- [ ] Created `pkg/authz/middleware.go`
- [ ] Moved authentication logic from `pkg/middleware/auth.go`
- [ ] Updated `pkg/api/routes.go` to use authz package
- [ ] Created `pkg/authz/rules.go` (optional)
- [ ] Removed `pkg/middleware/auth.go` (or kept for CORS)
- [ ] **Checkpoint**: Authentication still works

---

## Phase 7: Commanders ✅/❌

- [ ] Created `pkg/database/image_commander.go`
- [ ] Implemented `ImageCommander` interface
- [ ] Implemented `CreateImage()` command with transactions
- [ ] Updated handlers to use commanders
- [ ] Added more commands as needed (e.g., `UpdateImageStatus`)
- [ ] **Checkpoint**: Commanders handle transactions correctly

---

## Phase 8: Refactor Worker ✅/❌

- [ ] Updated `StartExtract()` signature
- [ ] Removed `*gorm.DB` parameter
- [ ] Added querier/commander parameters
- [ ] Replaced DB queries with repository calls
- [ ] Updated `ProcessImage()` function
- [ ] Used domain entity methods for state transitions
- [ ] Updated `cmd/main.go` to pass correct dependencies
- [ ] **Checkpoint**: Worker processes images using domain layer

---

## Phase 9: Directory Structure ✅/❌

- [ ] Created `cmd/server/` directory
- [ ] Moved `cmd/main.go` → `cmd/server/main.go`
- [ ] Updated build scripts/Dockerfiles
- [ ] Removed `pkg/models/models.go`
- [ ] Removed `pkg/middleware/auth.go` (if not needed)
- [ ] Moved test files to `test/rest/` (optional)
- [ ] Updated all imports
- [ ] **Checkpoint**: All imports resolve, app compiles

---

## Phase 10: Update Tests ✅/❌

- [ ] Created mock interfaces or used mocking library
- [ ] Updated handler tests to use mocks
- [ ] Updated worker tests to use mocks
- [ ] Created domain entity tests
- [ ] Created repository tests
- [ ] All tests pass
- [ ] **Checkpoint**: Test coverage maintained

---

## Final Verification

- [ ] All code compiles without errors
- [ ] All tests pass
- [ ] API endpoints work correctly
- [ ] Worker processes images correctly
- [ ] Authentication works
- [ ] No direct DB access in handlers or workers
- [ ] Domain layer has no external dependencies
- [ ] Clean architecture principles followed

---

## Notes

Use this space to note any issues, decisions, or deviations from the plan:
