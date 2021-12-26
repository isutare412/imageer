# Imageer API Server

## Generate OAS

```bash
# Install swag using go
make swag

# Generate OAS from source code
make docs
```

## Build

```bash
go build -o imageer_api cmd/main.go
```

## Run

```bash
# Run compiled executable
IMAGEER_CONFIG=configs/{target_profile}.yaml ./imageer_api

# ... or run without compile
IMAGEER_CONFIG=configs/{target_profile}.yaml go run cmd/main.go
```

## OpenAPI

- Swagger UI
  + `http://{api_server_url}/docs/index.html`
- OAS
  + `http://{api_server_url}/docs/swagger.json`
