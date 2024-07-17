# Use the official Golang image to create a build artifact.
FROM golang:1.16 AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and go sum files to /app directory
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source files from the current directory to the /app directory inside the container
COPY . .

# Build the Go app and create an executable named main
RUN go build -o main .

# Start a new stage from the official Golang image
FROM golang:1.16

# Set the working directory inside the container
WORKDIR /app

# Copy the pre-built binary file from the previous stage to /app
COPY --from=build /app/main /app/

# Copy the go.mod and go.sum files from the previous stage to /app
COPY --from=build /app/go.mod /app/go.sum /app/

# Copy the source files from the previous stage to /app
COPY --from=build /app/ /app/

# Copy the .env file to /app
COPY .env /app/

# Install necessary dependencies
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
