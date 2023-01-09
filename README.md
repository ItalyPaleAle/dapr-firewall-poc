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
  (cd app && ./run-daprd.sh)
  ```

2. Launch the app in another terminal window:  
  
  ```sh
  (cd app && ./run-app.sh)
  ```

You will see that Dapr will create an app channel and can send requests to the app, even though the app does not implement a server.

Next, try stopping the app and restarting it. And then, try stopping daprd and restarting it. The solution should recover automatically and quickly.

## Implementation details

At a high level, this is implemented by making the app (via the Dapr SDK) create an outbound TCP connection to the Dapr sidecar. Once the connection is up, the app starts a gRPC server on the established connection, and can begin accepting requests from the sidecar.

In details:

- The app is started without any gRPC or HTTP server, and `--app-port` is unset when starting `daprd`. However, there's a new flag for `daprd` called `--enable-callback-channel` which tells Dapr to expect the app to create a channel using the callback.
- There's a new method in the Dapr's runtime gRPC server called [`ConnectAppCallback`](https://github.com/ItalyPaleAle/dapr/blob/45a04142f826ce70d7bb290726da7e7be3cd4ec3/dapr/proto/runtime/v1/dapr.proto#L124). When the app starts, it creates a Dapr client (just as usual) which then invokes `ConnectAppCallback` on the sidecar.
  - When `ConnectAppCallback` is invoked, the sidecar [starts an ephemeral TCP listener](https://github.com/ItalyPaleAle/dapr/blob/45a04142f826ce70d7bb290726da7e7be3cd4ec3/pkg/grpc/api_connectappcallback.go#L31-L107), on a random port. It responds to the app's gRPC call with the port number.
  - The app then has a certain amount of time (currently, 10s) to establish a TCP connection to the ephemeral listener the sidecar has started. For the app, this is an outbound connection so it does not need any open firewall port (however, the sidecar needs to have the port open).
  - The port the sidecar opens is random by default. If a specific port needs to be used, daprd can be started with the `--callback-channel-port 1234` flag.
  - Once the TCP connection is established, Dapr automatically turns that into a "client connection" and creates a gRPC client on that.
  - Likewise, the app creates a gRPC server on the active TCP connection.
- Once the callback channel connection is established, Dapr invokes the `Ping` method on the app, which is a gRPC streaming call that is used to detect when the callback channel connection drops. This is necessary because if the connection drops, the app needs to have a way to detect that and re-connect to Dapr.
- All of the above are handled by the Dapr SDK automatically: the app just needs to invoke [`NewServiceFromCallbackChannel`](https://github.com/ItalyPaleAle/dapr-go-sdk/blob/e1ede39920d59860e183d9412796e5971183b0f1/service/grpc/service.go#L59-L87) and pass the existing client connection.

Here are the code diffs that make this possible:

- [dapr/dapr](https://github.com/ItalyPaleAle/dapr/compare/master...firewall)
- [dapr/go-sdk](https://github.com/ItalyPaleAle/dapr-go-sdk/compare/main...firewall)

Check out the demo app's [`main.go`](https://github.com/ItalyPaleAle/dapr-firewall-poc/blob/main/app/main.go) to see an example of how this is used.

## Current staus

- [X] POC working E2E
- SDK support:
  - [X] Go
  - [ ] .NET
  - [ ] Java
  - [ ] JavaScript - Possibly blocked due to grpc/grpc-node#2317
  - [ ] Python
  - [ ] Rust
- [X] Handle automatic reconnections if the connection drops
- [X] Unit tests
- [ ] E2E tests
