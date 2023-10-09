# Use the official Golang image as the base image
FROM golang:1.21.1-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy code to /app
COPY . /app

# Install the go sql driver
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/gorilla/mux 

RUN go build -o app

# Set go111module env variable
#ENV GO111MODULE=auto

CMD ["./app"]
