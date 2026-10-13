# Toolchain for `make check-docker`: Go plus libvips headers, same base as
# processor.Dockerfile, for hosts that cannot install libvips (e.g. Rocky 10).
FROM golang:1.26-trixie

RUN apt-get update && apt-get install -y --no-install-recommends \
    libvips-dev \
    && rm -rf /var/lib/apt/lists/*

# Caches live in the imageer-check-cache volume. A new volume copies these
# world-writable directories, so the container can run as the host user.
RUN mkdir -p /cache/home /cache/mod /cache/build /cache/lint && chmod -R 1777 /cache
