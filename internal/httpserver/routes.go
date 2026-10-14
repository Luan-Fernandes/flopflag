package httpserver

import (
	"net/http"

	"boliplate-go/internal/controller"
	"boliplate-go/internal/domain"
	"boliplate-go/internal/httpserver/middleware"
	"boliplate-go/internal/httpx"
	"boliplate-go/internal/security"

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
