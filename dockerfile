# Use the official Golang image as the base image
FROM golang:1.21.1-alpine

# Set the working directory inside the container
WORKDIR /app

# Install Git to clone the GitHub repository
RUN apk update && apk add --no-cache git

# Clone your Go application repository
RUN git clone https://github.com/obiMadu/hngX-stage2.git .

# Install the go sql driver
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/gorilla/mux 

RUN go build -o app

# Set go111module env variable
#ENV GO111MODULE=auto

CMD ["./app"]