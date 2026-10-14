package controller

import (
	"strconv"

	"github.com/Luan-Fernandes/flopflag/internal/httpserver/middleware"
	"github.com/Luan-Fernandes/flopflag/internal/httpx"
	"github.com/Luan-Fernandes/flopflag/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	users *service.UserService
}

func NewUserController(users *service.UserService) *UserController {
	return &UserController{users: users}
}

func (c *UserController) Me(ctx *gin.Context) {
	userID := middleware.UserID(ctx)
	if userID == uuid.Nil {
		httpx.Unauthorized(ctx, "não autenticado", nil)
		return
	}

	u, err := c.users.GetByID(ctx.Request.Context(), userID)
	if err != nil {
		httpx.FailErr(ctx, err)
		return
	}

	httpx.OK200(ctx, "", toUserView(*u))
}

func (c *UserController) List(ctx *gin.Context) {
	limit, offset := 50, 0

	if v := ctx.Query("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 || n > 200 {
			httpx.BadRequest(ctx, "limit inválido (1-200)", nil)
			return
		}
		limit = n
	}

	if v := ctx.Query("offset"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			httpx.BadRequest(ctx, "offset inválido", nil)
			return
		}
		offset = n
	}

	items, err := c.users.List(ctx.Request.Context(), limit, offset)
	if err != nil {
		httpx.FailErr(ctx, err)
		return
	}

	out := make([]any, 0, len(items))
	for _, u := range items {
		out = append(out, toUserView(u))
	}

	httpx.OK200(ctx, "", gin.H{
		"items": out,
		"meta": gin.H{
			"limit":  limit,
			"offset": offset,
		},
	})
}
