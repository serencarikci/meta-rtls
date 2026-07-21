package tenant

import "time"

type Tenant struct {
	ID           string    `json:"id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	ProfileScale string    `json:"profileScale"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
}

type CreateTenantRequest struct {
	Code         string `json:"code" binding:"required,min=2,max=64"`
	Name         string `json:"name" binding:"required,min=2,max=255"`
	ProfileScale string `json:"profileScale" binding:"required,oneof=SMALL MEDIUM LARGE"`
}
