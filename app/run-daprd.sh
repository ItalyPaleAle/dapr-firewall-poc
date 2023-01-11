#!/bin/sh

LOG_LEVEL=${1:-"debug"}
APP_NAME=firewallpoc

~/.dapr/bin/daprd \
    --app-id $APP_NAME \
    --dapr-http-port 3602 \
    --dapr-grpc-port 6602 \
    --metrics-port 9090 \
    --enable-callback-channel \
    --components-path ./components \
    --enable-app-health-check=true \
    --app-health-check-path=/healthz \
    --app-health-probe-interval=5 \
    --enable-profiling \
    --log-level "$LOG_LEVEL"
