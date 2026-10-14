package controller

import (
	"boliplate-go/internal/httpx"
	"boliplate-go/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	auth *service.AuthService
}

func NewAuthController(auth *service.AuthService) *AuthController {
	return &AuthController{auth: auth}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req registerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(ctx, "payload inválido", gin.H{"error": err.Error()})
		return
	}

	res, err := c.auth.Register(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		httpx.FailErr(ctx, err)
		return
	}

	httpx.Created(ctx, "", gin.H{
		"token": res.Token,
		"user":  toUserView(res.User),
	})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(ctx, "payload inválido", gin.H{"error": err.Error()})
		return
	}

	res, err := c.auth.Login(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		httpx.FailErr(ctx, err)
		return
	}

	httpx.OK200(ctx, "", gin.H{
		"token": res.Token,
		"user":  toUserView(res.User),
	})
}
