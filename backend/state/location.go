package state

import "time"

type Location struct {
	ID       int64
	Name     string
	Source   string
	Status   string
	Category string

	Latitude  float64
	Longitude float64
	Altitude  float64
	Speed     float64
	Course    float64
	Accuracy  float64
	Valid     bool

	Motion        *bool
	Distance      *float64
	TotalDistance *float64

	BatteryLevel *float64
	Charging     *bool

	ObservedAt time.Time
	ServerTime time.Time
	ReceivedAt time.Time
}
