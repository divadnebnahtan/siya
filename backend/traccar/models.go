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
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	UniqueID string `json:"uniqueId"`
	Status   string `json:"status"`
	Disabled bool   `json:"disabled"`

	LastUpdate *time.Time `json:"lastUpdate"`
	PositionID *int64     `json:"positionId"`
	GroupID    *int64     `json:"groupId"`

	Phone    *string `json:"phone"`
	Model    *string `json:"model"`
	Contact  *string `json:"contact"`
	Category *string `json:"category"`

	Attributes map[string]any `json:"attributes"`
}

type GetDevicesResponse []Device

type GetPositionsRequest struct {
	DeviceID *int64
	From     *time.Time
	To       *time.Time
	IDs      []int64
}

type Position struct {
	ID       int64  `json:"id"`
	DeviceID int64  `json:"deviceId"`
	Protocol string `json:"protocol"`

	DeviceTime time.Time `json:"deviceTime"`
	FixTime    time.Time `json:"fixTime"`
	ServerTime time.Time `json:"serverTime"`

	Valid     bool    `json:"valid"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
	Speed     float64 `json:"speed"`
	Course    float64 `json:"course"`
	Address   *string `json:"address"`
	Accuracy  float64 `json:"accuracy"`

	Network    map[string]any `json:"network"`
	Attributes map[string]any `json:"attributes"`
}

type GetPositionsResponse []Position
