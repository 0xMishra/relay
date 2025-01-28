#!/bin/bash

export GIT_REPO_URL="$GIT_REPO_URL"

git clone "$GIT_REPO_URL" /home/app/output

cd output

ls -la

cd ..

exec go run cmd/builder/main.go
