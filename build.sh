#!/usr/bin/env bash

rm riskclient
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o riskclient
