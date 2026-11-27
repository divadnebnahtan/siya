package state

import (
	"fmt"
	"math"
	"time"
)

type Location struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Source   string    `json:"source"`
	Status   string    `json:"status"`
	Category string    `json:"category,omitempty"`
	Position *Position `json:"position,omitempty"`
}

type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"` // m
	Speed     float64 `json:"speed"`    // km/h
	Course    float64 `json:"course"`   // deg
	Accuracy  float64 `json:"accuracy"` // m
	Valid     bool    `json:"valid"`

	Motion        *bool    `json:"motion,omitempty"`
	Distance      *float64 `json:"distance,omitempty"`      // m
	TotalDistance *float64 `json:"totalDistance,omitempty"` // m
	BatteryLevel  *float64 `json:"batteryLevel,omitempty"`  // %
	Charging      *bool    `json:"charging,omitempty"`

	ObservedAt time.Time `json:"observedAt"`
	ServerTime time.Time `json:"serverTime,omitzero"`
	ReceivedAt time.Time `json:"receivedAt"`
}

func (p Position) validate() error {
	for _, field := range []struct {
		name            string
		value, min, max float64
	}{
		{"latitude", p.Latitude, -90, 90},
		{"longitude", p.Longitude, -180, 180},
		{"altitude", p.Altitude, -math.MaxFloat64, math.MaxFloat64},
		{"speed", p.Speed, 0, math.MaxFloat64},
		{"course", p.Course, 0, 360},
		{"accuracy", p.Accuracy, 0, math.MaxFloat64},
	} {
		if !inRange(field.value, field.min, field.max) {
			return fmt.Errorf("invalid %s", field.name)
		}
	}

	for _, field := range []struct {
		name  string
		value *float64
		max   float64
	}{
		{"distance", p.Distance, math.MaxFloat64},
		{"total distance", p.TotalDistance, math.MaxFloat64},
		{"battery level", p.BatteryLevel, 100},
	} {
		if field.value != nil && !inRange(*field.value, 0, field.max) {
			return fmt.Errorf("invalid %s", field.name)
		}
	}
	if p.ObservedAt.IsZero() || p.ReceivedAt.IsZero() {
		return fmt.Errorf("position requires observation and receipt times")
	}
	for _, stamp := range []time.Time{p.ObservedAt, p.ServerTime, p.ReceivedAt} {
		year := stamp.UTC().Year()
		if year < 0 || year > 9999 {
			return fmt.Errorf("position timestamp is outside the JSON year range")
		}
	}
	return nil
}

func inRange(value, min, max float64) bool {
	return value >= min && value <= max
}

func (l Location) clone() Location {
	if l.Position != nil {
		p := *l.Position
		p.Motion = clonePointer(p.Motion)
		p.Distance = clonePointer(p.Distance)
		p.TotalDistance = clonePointer(p.TotalDistance)
		p.BatteryLevel = clonePointer(p.BatteryLevel)
		p.Charging = clonePointer(p.Charging)
		l.Position = &p
	}
	return l
}

func clonePointer[T any](p *T) *T {
	if p == nil {
		return nil
	}
	value := *p
	return &value
}
