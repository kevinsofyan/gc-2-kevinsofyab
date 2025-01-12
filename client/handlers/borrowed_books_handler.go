package handlers

import (
	"context"
	"net/http"
	"time"

	pb "gc-buku/proto"

	"github.com/labstack/echo/v4"
)

type BorrowedBooksHandler struct {
	grpcClient pb.BookServiceClient
}

func NewBorrowedBooksHandler(client pb.BookServiceClient) *BorrowedBooksHandler {
	return &BorrowedBooksHandler{grpcClient: client}
}

func (h *BorrowedBooksHandler) BorrowBook(c echo.Context) error {
	userID := c.Get("user_id").(string)
	bookID := c.Param("book_id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.grpcClient.BorrowBook(ctx, &pb.BorrowBookRequest{
		BorrowedBook: &pb.BorrowedBook{
			BookId:       bookID,
			UserId:       userID,
			BorrowedDate: time.Now().Format(time.RFC3339),
		},
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.BorrowedBook)
}

func (h *BorrowedBooksHandler) ReturnBook(c echo.Context) error {
	borrowID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.grpcClient.ReturnBook(ctx, &pb.ReturnBookRequest{
		Id: borrowID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.BorrowedBook)
}
