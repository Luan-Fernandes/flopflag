package httpserver

import (
	"net/http"

	"github.com/Luan-Fernandes/flopflag/internal/controller"
	"github.com/Luan-Fernandes/flopflag/internal/domain"
	"github.com/Luan-Fernandes/flopflag/internal/httpserver/middleware"
	"github.com/Luan-Fernandes/flopflag/internal/httpx"
	"github.com/Luan-Fernandes/flopflag/internal/security"

	"github.com/gin-gonic/gin"
)

func registerRoutes(
	r *gin.Engine,
	tokens *security.TokenManager,
	auth *controller.AuthController,
	users *controller.UserController,
) {
	r.GET("/health", func(c *gin.Context) {
		httpx.OK200(c, "ok", gin.H{"status": "ok"})
	})

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", auth.Register)
		authGroup.POST("/login", auth.Login)
	}

	protected := r.Group("/")
	protected.Use(middleware.AuthRequired(tokens))
	{
		protected.GET("/me", users.Me)

		admin := protected.Group("/admin")
		admin.Use(middleware.RequireRoles(domain.RoleAdmin))
		{
			admin.GET("/users", users.List)
			admin.GET("/health", func(c *gin.Context) {
				c.Status(http.StatusNoContent)
			})
		}
	}
}
