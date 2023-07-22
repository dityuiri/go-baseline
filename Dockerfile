FROM golang:1.20

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project directory into the container's working directory
COPY . .

# Copy the local.env file into the container
COPY config/local.env /app/config/local.env

# Ensure Go modules are enabled and download dependencies
RUN go mod download

# Build the application
RUN go build -o main .

# Expose ports for HTTP and gRPC
EXPOSE 8080 50051

CMD ["/app/main"]
