# Image Processor Server

## Build

```bash
go build -o imageer_processor cmd/main.go
```

## Run

```bash
# Run compiled executable
IMAGEER_CONFIG=configs/{target_profile}.yaml ./imageer_processor

# ... or run without compile
IMAGEER_CONFIG=configs/{target_profile}.yaml go run cmd/main.go
```
