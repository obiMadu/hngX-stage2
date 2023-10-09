# Use the official Golang image as the base image
FROM golang:1.21.1-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy code to /app
COPY . /app

# Build the app
RUN go build -o app

CMD ["./app"]
