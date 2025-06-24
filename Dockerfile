# syntax=docker/dockerfile:1

FROM golang:1.24.4-alpine3.22 AS build

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /cloudflare-updater

FROM alpine:3.22 AS final

COPY --from=build /cloudflare-updater /cloudflare-updater

RUN apk add --no-cache ca-certificates curl

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 5000

# healthcheck localhost:$DYNDNS_LISTEN_ADDR/ping (default 5000)
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s \
  CMD curl -f http://localhost:5000/ping || exit 1

# Run
CMD ["/cloudflare-updater"]