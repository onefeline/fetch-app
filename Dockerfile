# syntax=docker/dockerfile:1
FROM golang:1.21.0

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /fetch-app

EXPOSE 8080

# Run
CMD ["/fetch-app"]