FROM node:18-alpine AS builder

# Install protoc, git, and other dependencies
RUN apk add --no-cache protobuf git curl

# Download and install protoc-gen-grpc-web binary
RUN curl -L -o /usr/local/bin/protoc-gen-grpc-web https://github.com/grpc/grpc-web/releases/download/1.5.0/protoc-gen-grpc-web-1.5.0-linux-x86_64 && \
    chmod +x /usr/local/bin/protoc-gen-grpc-web

# Set working directory
WORKDIR /app

# Copy package.json and package-lock.json for dependency caching
COPY package.json package-lock.json ./

# Install dependencies
RUN npm install

# Copy the rest of the source code
COPY . .

# Generate gRPC-Web files
RUN chmod +x gen_front_proto.sh && ./gen_front_proto.sh

# Build the frontend (if you have a build step)
# RUN npm run build  # Uncomment if you have a build script

# Final image
FROM node:18-alpine

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Create a non-root user for security
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy built files from builder stage
COPY --from=builder /app .

# Switch to non-root user
USER appuser

# Expose port (assuming 3000 for the frontend dev server)
EXPOSE 3000

# Start the frontend
CMD ["npm", "start"]
