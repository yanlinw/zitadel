#!/bin/bash

set -e

CGO_ENABLED=0 go build \
  -a \
  -installsuffix cgo \
  -o backupctl \
  ./cmd/backupctl/main.go
