#!/bin/bash

set -e

./build/cr-backup/build.sh
exec ./backupctl "$@"
