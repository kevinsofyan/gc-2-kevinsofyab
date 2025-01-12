package handlers

import (
	"context"
	"net/http"
	"time"

	pb "gc-buku/proto"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	grpcClient pb.BookServiceClient
}

func NewUserHandler(client pb.BookServiceClient) *UserHandler {
	return &UserHandler{grpcClient: client}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	type CreateUserRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	req := new(CreateUserRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.grpcClient.CreateUser(ctx, &pb.CreateUserRequest{
		User: &pb.User{
			Username: req.Username,
			Password: req.Password,
		},
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.User)
}

func (h *UserHandler) Login(c echo.Context) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.grpcClient.Login(ctx, &pb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": resp.Token,
		"user":  resp.User,
	})
}

func (h *UserHandler) GetUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.grpcClient.GetUser(ctx, &pb.GetUserRequest{
		Id: c.Param("id"),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.User)
}
