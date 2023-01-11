#!/bin/sh

export DAPR_GRPC_PORT=6602
export APP_NAME=firewallpoc

go run .
