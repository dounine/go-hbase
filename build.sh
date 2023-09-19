#!/bin/bash
CGO_ENABLED=1 go build -ldflags '-linkmode "external" -extldflags "-static"' main.go