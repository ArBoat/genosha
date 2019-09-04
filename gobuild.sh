#!/bin/bash

export GIN_MODE=release
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o genosha .