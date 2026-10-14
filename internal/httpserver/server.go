package httpserver

import (
	"time"

	"boliplate-go/internal/config"
	"boliplate-go/internal/controller"
	"boliplate-go/internal/repository"
	"boliplate-go/internal/security"
	"boliplate-go/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	Router *gin.Engine
}

func New(cfg config.Config, db *gorm.DB) *Server {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	userRepo := repository.NewUserRepository(db)

	tokens := security.NewTokenManager(cfg.JWT.Secret, time.Duration(cfg.JWT.ExpiresMinutes)*time.Minute)

	authSvc := service.NewAuthService(userRepo, tokens)
	userSvc := service.NewUserService(userRepo)

	authController := controller.NewAuthController(authSvc)
	userController := controller.NewUserController(userSvc)

	registerRoutes(r, tokens, authController, userController)

	return &Server{Router: r}
}
