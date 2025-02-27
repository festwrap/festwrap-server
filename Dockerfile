ARG GO_VERSION="1.24.0-bookworm"
FROM golang:${GO_VERSION}

ARG PORT=8080
ARG USERNAME=server
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME

USER $USERNAME

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal

EXPOSE $PORT

CMD ["go", "run", "/app/cmd/main.go"]
