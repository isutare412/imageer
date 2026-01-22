# Build stage
FROM node:24-slim AS builder

WORKDIR /app

# Install pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

# Copy package files
COPY package.json pnpm-lock.yaml ./

# Install dependencies
RUN pnpm install --frozen-lockfile

# Copy source code and config files
COPY src/ ./src/
COPY static/ ./static/
COPY svelte.config.js tsconfig.json vite.config.ts ./

# Build the application
RUN pnpm run build

# Production stage
FROM node:24-slim AS production

WORKDIR /app

# Install pnpm for production dependencies
RUN corepack enable && corepack prepare pnpm@latest --activate

# Copy package files
COPY package.json pnpm-lock.yaml ./

# Install production dependencies only
RUN pnpm install --prod --frozen-lockfile

# Copy built application from builder stage
COPY --from=builder /app/build ./build

# Set environment variables
ENV NODE_ENV=production
ENV PORT=3000
ENV PUBLIC_API_BASE_URL=

EXPOSE 3000

# Run the application
CMD ["node", "build"]
