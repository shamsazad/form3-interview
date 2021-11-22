# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest as builder

# Add Maintainer Info
LABEL maintainer="Shams Azad"

# Set the Current Working Directory inside the container
WORKDIR /test-api

ADD . .
RUN go mod download

CMD "go" "test" "-v" "./..."