package main

import (
	"context"
	"fmt"

	"github.com/divadnebnahtan/siya/backend/config"
	"github.com/divadnebnahtan/siya/backend/traccar"
)

func main() {
	fmt.Println("Siya, world!")

	env := config.Load()

	traccarClient := traccar.NewClient(env.TraccarBaseUrl, env.GetAuthString())

	devices, err := traccarClient.GetDevices(context.Background(), traccar.GetDevicesRequest{})
	if err != nil {
		panic(err)
	}
	for _, device := range devices {
		fmt.Println(device)
	}

	fmt.Println()

	positions, err := traccarClient.GetPositions(context.Background(), traccar.GetPositionsRequest{})
	if err != nil {
		panic(err)
	}
	for _, position := range positions {
		fmt.Println(position, "\n")
	}
}
