package main

import (
	"fmt"
	"os"

	"github.com/divadnebnahtan/siya/backend/config"
	"github.com/divadnebnahtan/siya/backend/traccar"
)

func main() {
	if err := run(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(output *os.File) error {
	fmt.Fprintln(output, "Siya, world!")

	env, err := config.Load()
	if err != nil {
		return err
	}
	client, err := traccar.NewClient(env.TraccarBaseURL, env.GetAuthString())
	if err != nil {
		return fmt.Errorf("configure Traccar client: %w", err)
	}

	devices, err := client.GetDevices(traccar.GetDevicesRequest{})
	if err != nil {
		return err
	}
	for _, device := range devices {
		fmt.Fprintln(output, device)
	}
	fmt.Fprintln(output)

	positions, err := client.GetPositions(traccar.GetPositionsRequest{})
	if err != nil {
		return err
	}
	for _, position := range positions {
		fmt.Fprintf(output, "%v \n\n", position)
	}
	return nil
}
