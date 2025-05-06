FROM golang:1.24-alpine AS builder

# Allow overriding which cmd to build
ARG PROG_PATH=cmd/server/main.go

WORKDIR /app
RUN apk add --no-cache git ca-certificates

# Dependencies first, so Docker layer caching works nicely
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build \
  -trimpath \
  -o /app/bin/$(basename $(dirname ${PROG_PATH})) \
  "${PROG_PATH}"

# Final, minimal image
FROM scratch
WORKDIR /app

# copy certs in case you need HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# copy whatever binary was built
# e.g. /app/bin/telemetrygateway, /app/bin/telemetrygenerator, etc.
COPY --from=builder /app/bin/ /app/bin/

# default; override with `docker run ... /app/bin/client`
ENTRYPOINT ["/app/bin/telemetrygateway"]
