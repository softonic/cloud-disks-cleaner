FROM golang:1.18 as build

ENV GO111MODULE=on

WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /usr/local/bin/cloud-disks-cleaner .


FROM debian:stretch-slim
# Update stretch repositories
RUN sed -i -e 's/deb.debian.org/archive.debian.org/g' \
           -e 's|security.debian.org|archive.debian.org/|g' \
           -e '/stretch-updates/d' /etc/apt/sources.list
RUN apt-get -qqq update \
    && apt-get -qqq -y install ca-certificates \
    && update-ca-certificates \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*


COPY --from=build /usr/local/bin/cloud-disks-cleaner /usr/local/bin/cloud-disks-cleaner

ENTRYPOINT ["/usr/local/bin/cloud-disks-cleaner"]
