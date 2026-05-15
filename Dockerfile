ARG GO_VERSION=1.26.2

FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /src

COPY go.work ./
COPY apps ./apps
COPY packages ./packages

ARG APP_DIR
ARG APP_CMD
RUN test -n "${APP_DIR}" && test -n "${APP_CMD}"
RUN go build -C "${APP_DIR}" -o "/out/app" "${APP_CMD}"

FROM alpine:3.22
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /out/app /app/app
ENTRYPOINT ["/app/app"]
