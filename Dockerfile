# Dockerfile
# Stage 1: Build the application
FROM golang:1.22 AS builder1

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o plugin main.go

# Dockerfile
# Stage 2: Build privateer
FROM golang:1.22 AS builder2
WORKDIR /app
RUN git clone https://github.com/privateerproj/privateer
WORKDIR /app/privateer
RUN make go-build

# Dockerfile
# Stage 3: Build wireframe
FROM golang:1.22 AS builder3
WORKDIR /app
RUN git clone https://github.com/privateerproj/raid-wireframe
WORKDIR /app/raid-wireframe
RUN make go-build


# Stage 4: Create final image
FROM ubuntu:jammy-20240627.1
COPY --from=builder1 /app/plugin /compliance-framework/plugin
COPY --from=builder2 /app/privateer/privateer /usr/bin
COPY --from=builder3 /app/raid-wireframe/SVC /root/privateer/bin/svc


