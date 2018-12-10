#!/bin/bash -
for GOOS in darwin linux; do
    GOOS=$GOOS GOARCH=amd64 go build -o bin/fargate-deploy-$GOOS-amd64 cmd/fargate-deploy/*.go
done
