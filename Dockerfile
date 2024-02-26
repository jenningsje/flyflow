# Use the official Go image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.17 as builder

# Copy local code to the container image.
WORKDIR /go/src/github.com/flyflow-devs/flyflow
COPY . .

# Build the command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with 'go get'.)
RUN go build -v -o flyflow

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/flyflow-devs/flyflow/flyflow /flyflow

# Run the web service on container startup.
CMD ["/flyflow"]
