# Use the Go alpine base image
FROM golang:alpine

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project into the container
COPY . .

# Navigate to the directory where main.go is located
WORKDIR /app/cmd/gin

# Build the Go application and output it as a binary named
RUN go build -o /app/gocleanarch

# Ensure the binary has executable permissions
RUN chmod +x /app/gocleanarch

# Expose the port that the Go app will use
EXPOSE 8080

# Set the default command to run the Go application
CMD ["/app/gocleanarch"]