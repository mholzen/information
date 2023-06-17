# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.20.4

# Copy the local package files to the container's workspace.
ADD . /go/src/my-app

# Set the current working directory inside the container
WORKDIR /go/src/my-app

# Build the application inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go mod tidy
RUN go build -o main server/*

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the server command by default when the container starts.
ENTRYPOINT ["/go/src/my-app/main"]
