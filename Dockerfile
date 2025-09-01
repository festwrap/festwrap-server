ARG GO_VERSION="1.24.0-bookworm"
FROM golang:${GO_VERSION} AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal

# Do not use CGO for statically linked self-contained binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd

# https://labs.iximiuz.com/tutorials/gcr-distroless-container-images
FROM gcr.io/distroless/static:nonroot

ARG PORT=8080
EXPOSE $PORT

COPY --from=builder /app/server /server

ENTRYPOINT ["/server"]
