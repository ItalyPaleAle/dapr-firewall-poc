# Dapr POC: apps behind a firewall

This is a POC of adding support to [Dapr](https://dapr.io) for using apps that are behind a firewall–or, more generally, do not have a gRPC or HTTP server running.

The goal is to solve the problem highlighted in dapr/dapr#5392

## Running the POC

> These steps are optimized for running on Linux or macOS.

### Set up the environment

First, clone this repo:

```sh
git clone https://github.com/ItalyPaleAle/dapr-firewall-poc dapr-firewall-poc
cd dapr-firewall-poc
```

Clone Dapr from my fork:

```sh
git clone https://github.com/ItalyPaleAle/dapr dapr
(cd dapr && git checkout firewall)
```

Clone the Dapr Go SDK from my fork:

```sh
git clone https://github.com/ItalyPaleAle/dapr-go-sdk go-sdk
(cd go-sdk && git checkout firewall)
```

> You need to have the Dapr CLI installed and have already executed `dapr init`

Build Dapr and "install" the binary:

```sh
./build-dapr.sh
```

### Run the app

You will need 2 terminal windows to launch Dapr and the app.

1. Launch Dapr in one terminal window:  
  
  ```sh
  cd app && ./run-daprd.sh
  ```

2. Launch the app in another terminal window:  
  
  ```sh
  cd app && ./run-app.sh
  ```

## Implementation details

## Current staus

- [X] POC working E2E
- SDK support:
  - [X] Go
  - [ ] .NET
  - [ ] Java
  - [ ] JavaScript
  - [ ] Python
  - [ ] Rust
- [ ] Handle automatic reconnections if the connection drops
- [ ] Unit tests
- [ ] E2E tests
