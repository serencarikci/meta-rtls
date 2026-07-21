package identity

import "time"

type User struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenantId"`
	Email        string    `json:"email"`
	DisplayName  string    `json:"displayName"`
	Role         string    `json:"role"`
	Status       string    `json:"status"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
}

type LoginRequest struct {
	TenantCode string `json:"tenantCode" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
	AccessToken string    `json:"accessToken"`
	ExpiresAt   time.Time `json:"expiresAt"`
	User        User      `json:"user"`
}

type MeResponse struct {
	User       User   `json:"user"`
	TenantCode string `json:"tenantCode"`
	TenantName string `json:"tenantName"`
}
