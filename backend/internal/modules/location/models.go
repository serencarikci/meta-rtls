package location

import "time"

type PositionEvent struct {
	TenantID  string    `json:"tenantId"`
	TagID     string    `json:"tagId"`
	TagCode   string    `json:"tagCode"`
	FloorID   string    `json:"floorId"`
	X         float64   `json:"x"`
	Y         float64   `json:"y"`
	Timestamp time.Time `json:"timestamp"`
}

type LivePosition struct {
	TenantID string    `json:"tenantId"`
	TagID    string    `json:"tagId"`
	TagCode  string    `json:"tagCode"`
	FloorID  string    `json:"floorId"`
	X        float64   `json:"x"`
	Y        float64   `json:"y"`
	ZoneCode string    `json:"zoneCode,omitempty"`
	Updated  time.Time `json:"updatedAt"`
}

type ZoneBox struct {
	ID   string
	Code string
	MinX float64
	MinY float64
	MaxX float64
	MaxY float64
}

type SimTag struct {
	TagID   string
	TagCode string
	FloorID string
	X       float64
	Y       float64
	DirX    float64
	DirY    float64
}
