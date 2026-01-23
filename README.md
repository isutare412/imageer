# Imageer

Imageer is a self-hosted image processing service that automatically transforms
uploaded images into multiple variants. It uses presigned URLs for secure
uploads to S3 and processes images with [bimg](https://github.com/h2non/bimg)
(libvips wrapper) workers running on Kubernetes.

## Architecture

The system runs on Kubernetes with a gateway handling upload URLs and image
metadata, processors for image transformation, and Valkey for request streaming.
PostgreSQL stores image state and Google OIDC provides authentication.

![Deployment diagram](./docs/assets/deployment.drawio.png)

## Screenshots

### Image Management

Browse and manage images across projects. View processing status, variants, and
preview transformed images.

![Image management](./docs/assets/image-management.png)

### Project Presets

Define image processing presets per project. Configure output format, quality,
dimensions, and fit mode for each variant.

![Project presets](./docs/assets/project-presets.png)

### Service Accounts

Manage API credentials for external services. Control access scope per project
with expiration settings.

![Service accounts](./docs/assets/service-accounts.png)

### Upload Tester

Test the upload flow with presigned URLs. Select a project, upload an image, and
verify the processing pipeline.

![Upload tester](./docs/assets/upload-tester.png)

## Development

```bash
make help
```
