#!/bin/sh

LOG_LEVEL=${1:-"debug"}
APP_NAME=firewallpoc

docker run \
  --rm \
  -p 3602:3602 \
  -p 6602:6602 \
  -p 2002:2002 \
  dapr-firewall:latest\
    /daprd \
      --app-id $APP_NAME \
      --dapr-http-port 3602 \
      --dapr-grpc-port 6602 \
      --metrics-port 9090 \
      --enable-callback-channel \
      --callback-channel-port=2002 \
      --components-path ./components \
      --enable-app-health-check=true \
      --app-health-check-path=/healthz \
      --app-health-probe-interval=5 \
      --enable-profiling \
      --log-level "$LOG_LEVEL"
