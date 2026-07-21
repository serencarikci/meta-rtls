package rtlsconfig

import "time"

type Site struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenantId"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Timezone  string    `json:"timezone"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

type Building struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
	SiteID   string `json:"siteId"`
	Code     string `json:"code"`
	Name     string `json:"name"`
}

type Floor struct {
	ID         string  `json:"id"`
	TenantID   string  `json:"tenantId"`
	BuildingID string  `json:"buildingId"`
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	LevelIndex int     `json:"levelIndex"`
	WidthM     float64 `json:"widthM"`
	HeightM    float64 `json:"heightM"`
}

type Zone struct {
	ID       string  `json:"id"`
	TenantID string  `json:"tenantId"`
	FloorID  string  `json:"floorId"`
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	ZoneType string  `json:"zoneType"`
	MinX     float64 `json:"minX"`
	MinY     float64 `json:"minY"`
	MaxX     float64 `json:"maxX"`
	MaxY     float64 `json:"maxY"`
}

type CreateSiteRequest struct {
	Code     string `json:"code" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Timezone string `json:"timezone"`
}
