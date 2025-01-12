package handlers

import (
	"context"
	"net/http"
	"time"

	pb "gc-buku/proto"

	"github.com/labstack/echo/v4"
)

type BookHandler struct {
	grpcClient pb.BookServiceClient
}

func NewBookHandler(client pb.BookServiceClient) *BookHandler {
	return &BookHandler{grpcClient: client}
}

func (h *BookHandler) CreateBook(c echo.Context) error {
	type CreateBookRequest struct {
		Title         string `json:"title"`
		Author        string `json:"author"`
		PublishedDate string `json:"published_date"`
	}
	req := new(CreateBookRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.grpcClient.CreateBook(ctx, &pb.CreateBookRequest{
		Book: &pb.Book{
			Title:         req.Title,
			Author:        req.Author,
			PublishedDate: req.PublishedDate,
		},
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, resp.Book)
}

func (h *BookHandler) GetBook(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.grpcClient.GetBook(ctx, &pb.GetBookRequest{Id: c.Param("id")})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.Book)
}

func (h *BookHandler) UpdateBook(c echo.Context) error {
	type UpdateBookRequest struct {
		Title         string `json:"title"`
		Author        string `json:"author"`
		PublishedDate string `json:"published_date"`
		Status        string `json:"status"`
	}

	req := new(UpdateBookRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.grpcClient.UpdateBook(ctx, &pb.UpdateBookRequest{
		Book: &pb.Book{
			Id:            c.Param("id"),
			Title:         req.Title,
			Author:        req.Author,
			PublishedDate: req.PublishedDate,
			Status:        req.Status,
		},
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.Book)
}

func (h *BookHandler) DeleteBook(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.grpcClient.DeleteBook(ctx, &pb.DeleteBookRequest{
		Id: c.Param("id"),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"id": resp.Id})
}
