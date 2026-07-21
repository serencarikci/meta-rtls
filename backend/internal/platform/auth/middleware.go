package auth

import (
	"net/http"
	"strings"

	"github.com/denizyetis/meta-rtls/internal/platform/response"
	"github.com/gin-gonic/gin"
)

const (
	CTX_USER_ID   = "authUserID"
	CTX_TENANT_ID = "authTenantID"
	CTX_ROLE      = "authRole"
	CTX_EMAIL     = "authEmail"
)

func Middleware(tokens *TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Fail(c, http.StatusUnauthorized, "missing bearer token")
			return
		}
		claims, err := tokens.Parse(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, "invalid token")
			return
		}
		c.Set(CTX_USER_ID, claims.UserID)
		c.Set(CTX_TENANT_ID, claims.TenantID)
		c.Set(CTX_ROLE, claims.Role)
		c.Set(CTX_EMAIL, claims.Email)
		c.Next()
	}
}

func TenantID(c *gin.Context) string {
	v, _ := c.Get(CTX_TENANT_ID)
	s, _ := v.(string)
	return s
}

func UserID(c *gin.Context) string {
	v, _ := c.Get(CTX_USER_ID)
	s, _ := v.(string)
	return s
}
