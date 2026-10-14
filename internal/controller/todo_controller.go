package controller

import (
	"boliplate-go/internal/httpserver/middleware"
	"boliplate-go/internal/httpx"
	"boliplate-go/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TodoController struct {
	todos *service.TodoService
}

func NewTodoController(todos *service.TodoService) *TodoController {
	return &TodoController{todos: todos}
}

type createTodoRequest struct {
	Title string `json:"title" binding:"required"`
}

type updateTodoRequest struct {
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
}

func (c *TodoController) List(ctx *gin.Context) {
	userID := middleware.UserID(ctx)
	items, err := c.todos.List(ctx.Request.Context(), userID)
	if err != nil {
		httpx.FailErr(ctx, err)
		return
	}
	httpx.OK200(ctx, "", items)
}

func (c *TodoController) Get(ctx *gin.Context) {
	userID := middleware.UserID(ctx)
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		httpx.BadRequest(ctx, "id inválido", nil)
		return
	}

	t, err := c.todos.Get(ctx.Request.Context(), id, userID)
	if err != nil {
		httpx.FailErr(ctx, err)
		return
	}
	httpx.OK200(ctx, "", t)
}

func (c *TodoController) Create(ctx *gin.Context) {
	userID := middleware.UserID(ctx)

	var req createTodoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(ctx, "payload inválido", gin.H{"error": err.Error()})
		return
	}

	t, err := c.todos.Create(ctx.Request.Context(), userID, req.Title)
	if err != nil {
		httpx.FailErr(ctx, err)
		return
	}

	httpx.Created(ctx, "", t)
}

func (c *TodoController) Update(ctx *gin.Context) {
	userID := middleware.UserID(ctx)
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		httpx.BadRequest(ctx, "id inválido", nil)
		return
	}

	var req updateTodoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(ctx, "payload inválido", gin.H{"error": err.Error()})
		return
	}

	if req.Title == nil && req.Done == nil {
		httpx.BadRequest(ctx, "nada para atualizar", nil)
		return
	}

	t, err := c.todos.Update(ctx.Request.Context(), id, userID, req.Title, req.Done)
	if err != nil {
		httpx.FailErr(ctx, err)
		return
	}

	httpx.OK200(ctx, "", t)
}

func (c *TodoController) Delete(ctx *gin.Context) {
	userID := middleware.UserID(ctx)
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		httpx.BadRequest(ctx, "id inválido", nil)
		return
	}

	if err := c.todos.Delete(ctx.Request.Context(), id, userID); err != nil {
		httpx.FailErr(ctx, err)
		return
	}
	httpx.NoContent(ctx)
}
