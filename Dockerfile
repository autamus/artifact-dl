# Start from the latest golang base image
FROM ghcr.io/autamus/go:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Add Spack Built Packages to PATH
ENV PATH=/opt/view/bin:$PATH

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o artifact-dl .

# Start again with minimal envoirnment.
FROM ubuntu:latest

RUN apt-get update && \
    apt-get install -y ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /app

COPY --from=builder /app/artifact-dl /app/artifact-dl

# Command to run the executable
ENTRYPOINT ["/app/artifact-dl"]