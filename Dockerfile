FROM golang:1.18 as build

ENV GO111MODULE=on


WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /usr/local/bin/diskCleanupManager .


FROM debian:stretch-slim
RUN apt-get -qqq update \
    && apt-get -qqq -y install ca-certificates \
    && update-ca-certificates \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*


COPY --from=build /usr/local/bin/diskCleanupManager /usr/local/bin/diskCleanupManager

ENTRYPOINT ["/usr/local/bin/diskCleanupManager"]
