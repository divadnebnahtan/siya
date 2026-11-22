package main

import (
	"fmt"

	"github.com/divadnebnahtan/siya/backend/config"
	"github.com/divadnebnahtan/siya/backend/traccar"
)

func main() {
	fmt.Println("Siya, world!")

	env := config.Load()

	traccarClient := traccar.NewClient(env.TraccarBaseUrl, env.GetAuthString())

	// devices := getDevices(client, baseUrl, authStr)
	// for _, device := range devices {
	// 	fmt.Println("DEVICE")
	// 	for k, v := range device {
	// 		fmt.Println("-", k, ":", v)
	// 	}
	// }

	positions := traccarClient.GetPositions()
	for _, position := range positions {
		fmt.Println("POSITION")
		for k, v := range position {
			fmt.Println("-", k, ":", v)
		}
	}
}
