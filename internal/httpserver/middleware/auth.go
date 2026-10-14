package middleware

import (
	"strings"

	"github.com/Luan-Fernandes/flopflag/internal/httpx"
	"github.com/Luan-Fernandes/flopflag/internal/security"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ctxUserID = "auth.userID"
	ctxRole   = "auth.role"
)

func AuthRequired(tokens *security.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" {
			httpx.Unauthorized(c, "token ausente", nil)
			c.Abort()
			return
		}

		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
			httpx.Unauthorized(c, "token inválido", nil)
			c.Abort()
			return
		}

		userID, role, err := tokens.Parse(parts[1])
		if err != nil {
			httpx.Unauthorized(c, "token inválido", nil)
			c.Abort()
			return
		}

		c.Set(ctxUserID, userID)
		c.Set(ctxRole, role)
		c.Next()
	}
}

func RequireRoles(allowed ...string) gin.HandlerFunc {
	allowedSet := map[string]struct{}{}
	for _, r := range allowed {
		allowedSet[r] = struct{}{}
	}

	return func(c *gin.Context) {
		role, _ := c.Get(ctxRole)
		roleStr, _ := role.(string)
		if _, ok := allowedSet[roleStr]; !ok {
			httpx.Forbidden(c, "sem permissão", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

func UserID(c *gin.Context) uuid.UUID {
	v, ok := c.Get(ctxUserID)
	if !ok {
		return uuid.Nil
	}
	id, _ := v.(uuid.UUID)
	return id
}

func Role(c *gin.Context) string {
	v, ok := c.Get(ctxRole)
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}
