# Copilot Instructions for Imageer

## Architecture Overview
Imageer is a microservices architecture (MSA) project with two main services:
- **Gateway** (`internal/gateway/`): API gateway for managing image resources, projects, users
- **Processor** (`internal/processor/`): Event-driven image processor consuming S3 events

Services communicate via generated HTTP clients from OpenAPI specs. Shared types live in `pkg/` directory.

## Critical Code Generation Workflow
**⚠️ NEVER edit `.gen.go` files directly** - they are auto-generated and will be overwritten.

For any API changes:
1. Edit `internal/gateway/web/openapi.yaml` (single source of truth)
2. Run `make generate` to regenerate:
   - `internal/gateway/web/server.gen.go` (Echo handlers)
   - `pkg/gateway/client.gen.go` (HTTP client used by Processor)
   - `pkg/gateway/client.gen_mock.go` (test mocks)

## Project Structure Patterns
```
internal/{service}/     # Private service implementation
pkg/                   # Shared packages between services
  gateway/            # Generated client code (DO NOT EDIT)
  images/             # Image states and types
  users/              # User authorities
```

## Key Development Commands
```bash
make help              # Show available targets
make generate          # Generate all code from OpenAPI (REQUIRED after API changes)
```

## Image Processing Workflow
1. Client requests presigned URL from Gateway (`/api/v1/projects/{id}/images/upload-urls`)
2. Client uploads to S3 using presigned URL
3. S3 sends event to SQS queue
4. Processor service consumes SQS messages
5. Processor fetches transformation presets from Gateway API
6. Processor applies transformations, updates image state: `WAITING_UPLOAD` → `PROCESSING` → `READY`

## OpenAPI Conventions
- Use OpenAPI 3.1.1 syntax in `internal/gateway/web/openapi.yaml`
- Include Go-specific extensions: `x-go-type`, `x-go-type-import`
- For optional fields: Use `x-go-type-skip-optional-pointer: true` to avoid pointer types
- Echo server generation via `generate-server.yaml`
- Client generation via `generate-client.yaml` (consumed by Processor)

## Service Communication Patterns
- Processor service imports and uses client from `pkg/gateway/client.gen.go`
- Cross-service types defined in `pkg/` (e.g., `images.State`, `users.Authority`)
- Event-driven architecture: S3 → SQS → Processor → Gateway API calls

## Domain Models Location
- Core business logic in `internal/gateway/domain/` (image.go, project.go, user.go, service_account.go)
- Shared enums in `pkg/`: `images.State` (WAITING_UPLOAD, PROCESSING, READY), `users.Authority` (ADMIN, GUEST)
- Generated API models in `.gen.go` files (do not edit)

## Admin Operations
Batch reprocessing available via:
- `POST /api/v1/admin/projects/{projectId}/images/reprocess`
- Request body: `{"imageIds": ["uuid1", "uuid2", ...]}` or `{"reprocessAll": true}`

## Testing Approach
- Use generated mocks in `pkg/gateway/client.gen_mock.go`
- Test each service module independently
- Integration tests should cover S3 event → Processor → Gateway API workflow

## Code Generation Dependencies
- `github.com/oapi-codegen/oapi-codegen/v2` for OpenAPI code generation
- `go.uber.org/mock/mockgen` for mock generation
- Echo v4 framework for HTTP handling
