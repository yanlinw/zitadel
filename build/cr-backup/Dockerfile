#######################
## By default we build the prod enviroment
ARG ENV=prod

#######################
## Go base build
## Speed up this step by mounting your local go mod pkg directory
#######################
FROM golang:1.16 as go-base

WORKDIR src/github.com/caos/zitadel/
COPY go.mod go.sum ./
RUN go mod download

## Go prod build
FROM go-base as prod-go-build
COPY . .
RUN CGO_ENABLED=0 go build \
      -a \
      -installsuffix cgo \
      -o backupctl \
      ./cmd/backupctl/main.go

## Go dev build
FROM go-base as dev-go-build
RUN go get github.com/go-delve/delve/cmd/dlv

#######################
## Build binary
#######################
FROM alpine:latest as artifact

COPY --from=prod-go-build /go/src/github.com/caos/zitadel/backupctl /app/backupctl
RUN chmod a+x /app/backupctl

#######################
## Final Production Image
#######################
FROM alpine:latest as final
RUN apk update
RUN apk upgrade
RUN apk --no-cache add openssl tzdata tar bash unzip

ENV GLIBC_REPO=https://github.com/sgerrand/alpine-pkg-glibc
ENV GLIBC_VERSION=2.30-r0

# Compatibility with glibc which is used by cockroach linking
RUN set -ex && \
    apk --update add libstdc++ curl ca-certificates && \
    for pkg in glibc-${GLIBC_VERSION} glibc-bin-${GLIBC_VERSION}; \
        do curl -sSL ${GLIBC_REPO}/releases/download/${GLIBC_VERSION}/${pkg}.apk -o /tmp/${pkg}.apk; done && \
    apk add --allow-untrusted /tmp/*.apk && \
    rm -v /tmp/*.apk && \
    /usr/glibc-compat/sbin/ldconfig /lib /usr/glibc-compat/lib

# cockroach binary used to backup database
RUN mkdir /usr/local/lib/cockroach
RUN wget -qO- https://binaries.cockroachdb.com/cockroach-v21.1.7.linux-amd64.tgz \
    | tar  xvz && cp -i cockroach-v21.1.7.linux-amd64/cockroach /usr/local/bin/
RUN cp -i cockroach-v21.1.7.linux-amd64/lib/libgeos.so /usr/local/lib/cockroach/
RUN cp -i cockroach-v21.1.7.linux-amd64/lib/libgeos_c.so /usr/local/lib/cockroach/
RUN rm -r cockroach-v21.1.7.linux-amd64

RUN wget -q https://downloads.rclone.org/v1.35/rclone-v1.35-linux-amd64.zip -O ./rclone.zip
RUN unzip ./rclone.zip && cp ./rclone-v1.35-linux-amd64/rclone /usr/bin/rclone
RUN rm -r ./rclone.zip
RUN rm -r ./rclone-v1.35-linux-amd64
RUN chmod +x /usr/bin/rclone

COPY --from=artifact /app /
HEALTHCHECK NONE
ENTRYPOINT ["/backupctl"]
