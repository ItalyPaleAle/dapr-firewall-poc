package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sync/atomic"
	"time"

	daprClient "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
)

const failureRate float64 = 0.0

var (
	//plan    = []bool{true, false, false, false, false, true, true}
	plan = []bool{}
)

func main() {
	rand.Seed(time.Now().UnixMicro())

	// Create a Dapr client that establishes a gRPC connection
	fmt.Println("Connecting Dapr client")
	client, err := daprClient.NewClient()
	if err != nil {
		panic(err)
	}

	// In background, invoke this very same app every 2 seconds
	go func() {
		for {
			time.Sleep(2 * time.Second)
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			res, err := client.InvokeMethod(ctx, os.Getenv("APP_NAME"), "foo", "GET")
			cancel()
			fmt.Println("Invoke:", err, string(res))
		}
	}()

	// Create the service using a traditional listener
	// srv, err := daprd.NewService("127.0.0.1:" + os.Getenv("DAPR_GRPC_PORT"))

	// Create the service using the callback channel
	srv, err := daprd.NewServiceFromCallbackChannel(client)
	if err != nil {
		panic(err)
	}

	// Handler for the cron input binding message
	srv.AddBindingInvocationHandler("schedule", func(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
		fmt.Println("Received scheduled message")
		return nil, nil
	})

	// Handler for "foo" service invocation
	srv.AddServiceInvocationHandler("foo", func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
		fmt.Println("Received foo request")
		return &common.Content{
			Data: []byte("hello world"),
		}, nil
	})

	// Handler for health checks
	srv.AddHealthCheckHandler("/healthz", func(ctx context.Context) error {
		err := doHealthCheck()
		if err != nil {
			return err
		}
		return nil
	})

	// Start the gRPC server
	// This is a blocking call
	fmt.Println("Starting server")
	err = srv.Start()
	if err != nil {
		panic(err)
	}
}

var count atomic.Int64

func doHealthCheck() error {
	success := true
	v := count.Add(1)
	if v <= int64(len(plan)) {
		success = plan[v-1]
	} else {
		success = rand.Float64() > failureRate
	}

	if success {
		fmt.Println("Responding to health check request with success")
		return nil
	} else {
		fmt.Println("Responding to health check request with failure")
		return errors.New("simulated failure")
	}
}
