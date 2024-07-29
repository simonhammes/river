#!/usr/bin/env bash

set -euo pipefail

# TODO: Check if river is installed
# go install github.com/riverqueue/river/cmd/river@latest

# TODO: Use absolute path
source .env

river migrate-up --database-url "$DATABASE_URL"
