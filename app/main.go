package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
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

	fmt.Println("Connecting Dapr client")
	client, err := daprClient.NewClient()
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			time.Sleep(2 * time.Second)
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			res, err := client.InvokeMethod(ctx, "firewallpoc", "foo", "GET")
			cancel()
			fmt.Println("Invoke:", err, string(res))
		}
	}()

	srv, err := daprd.NewServiceFromCallbackChannel(client)
	if err != nil {
		panic(err)
	}

	srv.AddBindingInvocationHandler("schedule", func(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
		fmt.Println("Received scheduled message")
		return nil, nil
	})

	srv.AddServiceInvocationHandler("foo", func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
		fmt.Println("Received foo request")
		return &common.Content{
			Data: []byte("hello world"),
		}, nil
	})

	srv.AddHealthCheckHandler("/healthz", func(ctx context.Context) error {
		err := doHealthCheck()
		if err != nil {
			return err
		}
		return nil
	})

	// Blocking call
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
