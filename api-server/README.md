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
make build
```

## Run

```bash
make run-dev
```

## OpenAPI

- Swagger UI
  + `http://{api_server_url}/docs/index.html`
- OAS
  + `http://{api_server_url}/docs/swagger.json`
