# -----------------------------------------------------------------------------
# Build stage - Intermidiate  container
# -----------------------------------------------------------------------------
FROM golang:1.16-bullseye AS builder

# Setup working directory
WORKDIR /app

# Copy go.mod and go.sum files in containers /app folder
COPY go.mod go.sum ./

# Downloads and verify the packages in go.mod file
RUN go mod download
RUN go mod verify
# Copy source directory in containers /app folder
COPY . .

# Build the binary
RUN go build -o main main.go

# -----------------------------------------------------------------------------
# Build main Docker image
# -----------------------------------------------------------------------------
FROM golang:1.16-bullseye

# Setup working directory
WORKDIR /app

# Copy the file from build stage to run stage
COPY --from=builder /app/main /app/main
COPY migrate.linux-amd64 ./migrate
COPY db/migration ./migration
COPY start.sh ./start.sh


EXPOSE ${PORT}
CMD [ "./start.sh" ]





