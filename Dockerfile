# Using the official golang image as the base image
FROM golang:1.22.0-alpine AS builder

# Set up the working directory inside the container
WORKDIR /app

# Copy all files from the local directory to the working directory
COPY ./app .

# Also copy the wait-for-db.sh script and make it executable
COPY wait-for-db.sh /app/wait-for-db.sh

# Build the application
RUN go build -o /bin/main ./cmd/main.go

# Use a small image to run the app
FROM alpine:3.9

RUN apk add --no-cache bash

# Copy necessary files
COPY --from=builder /app/wait-for-db.sh /app/wait-for-db.sh
COPY ./app/config/config.json /config/config.json
COPY --from=builder /bin/main /bin/main

# Make the wait-for-db.sh script executable
RUN chmod +x /app/wait-for-db.sh
