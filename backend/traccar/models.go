package traccar

import "time"

type GetDevicesRequest struct {
	All               *bool
	UserID            *int64
	IDs               []int64
	UniqueIDs         []string
	ExcludeAttributes *bool
	Limit             *int
	Offset            *int
	Keyword           string
}

type Device struct {
	ID       int64
	Name     string
	UniqueID string
	Status   string
	Disabled bool

	LastUpdate *time.Time
	PositionID *int64
	GroupID    *int64

	Phone    *string
	Model    *string
	Contact  *string
	Category *string

	Attributes map[string]any
}

type GetDevicesResponse []Device

type GetPositionsRequest struct {
	DeviceID *int64
	From     *time.Time
	To       *time.Time
	IDs      []int64
}

type Position struct {
	ID       int64
	DeviceID int64
	Protocol string

	DeviceTime time.Time
	FixTime    time.Time
	ServerTime time.Time

	Valid     bool
	Latitude  float64
	Longitude float64
	Altitude  float64
	Speed     float64
	Course    float64
	Address   *string
	Accuracy  float64

	Network    map[string]any
	Attributes map[string]any
}

type GetPositionsResponse []Position
