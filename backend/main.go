package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/divadnebnahtan/siya/backend/config"
	"github.com/divadnebnahtan/siya/backend/state"
	"github.com/divadnebnahtan/siya/backend/traccar"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	env, err := config.Load()
	if err != nil {
		return err
	}
	client, err := traccar.NewClient(env.TraccarBaseURL, env.GetAuthString())
	if err != nil {
		return fmt.Errorf("configure Traccar client: %w", err)
	}

	var store state.Store
	err = refreshLocations(client, &store)
	if err != nil {
		return err
	}
	locations, _ := store.Snapshot()
	return json.NewEncoder(os.Stdout).Encode(locations)
}

func refreshLocations(client *traccar.Client, store *state.Store) error {
	devices, err := client.GetDevices(traccar.GetDevicesRequest{})
	if err != nil {
		return err
	}
	positions, err := client.GetPositions(traccar.GetPositionsRequest{})
	if err != nil {
		return err
	}
	receivedAt := time.Now()

	byID := make(map[int64]state.Location, len(devices))
	batch := make([]state.Location, 0, len(devices)+len(positions))
	for _, device := range devices {
		location := state.Location{
			ID: device.ID, Name: device.Name, Source: "traccar", Status: device.Status,
		}
		if device.Category != nil {
			location.Category = *device.Category
		}
		byID[device.ID] = location
		batch = append(batch, location)
	}
	for _, position := range positions {
		location, ok := byID[position.DeviceID]
		if !ok {
			// The device list can change between the two requests.
			continue
		}
		location.Position = &state.Position{
			Latitude: position.Latitude, Longitude: position.Longitude,
			Altitude: position.Altitude, Speed: position.Speed,
			Course: position.Course, Accuracy: position.Accuracy, Valid: position.Valid,
			Motion:        attribute[bool](position.Attributes, "motion"),
			Distance:      attribute[float64](position.Attributes, "distance"),
			TotalDistance: attribute[float64](position.Attributes, "totalDistance"),
			BatteryLevel:  attribute[float64](position.Attributes, "batteryLevel"),
			Charging:      attribute[bool](position.Attributes, "charge"),
			ObservedAt:    position.FixTime, ServerTime: position.ServerTime, ReceivedAt: receivedAt,
		}
		batch = append(batch, location)
	}
	err = store.Upsert(batch...)
	if err != nil {
		return fmt.Errorf("update locations: %w", err)
	}
	return nil
}

func attribute[T bool | float64](attributes map[string]any, name string) *T {
	value, ok := attributes[name].(T)
	if !ok {
		return nil
	}
	return &value
}
