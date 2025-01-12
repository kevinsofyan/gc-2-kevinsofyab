package routes

import (
	"gc-buku/client/handlers"
	"gc-buku/client/middleware"

	pb "gc-buku/proto"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, client pb.BookServiceClient) {
	// Handlers
	userHandler := handlers.NewUserHandler(client)
	bookHandler := handlers.NewBookHandler(client)
	borrowedBooksHandler := handlers.NewBorrowedBooksHandler(client)

	// Public routes
	e.POST("/register", userHandler.CreateUser)
	e.POST("/login", userHandler.Login)

	// Protected book routes
	books := e.Group("/books", middleware.Auth)
	{
		books.POST("", bookHandler.CreateBook)
		books.GET("/:id", bookHandler.GetBook)
		books.PUT("/:id", bookHandler.UpdateBook)
		books.DELETE("/:id", bookHandler.DeleteBook)
	}

	// Protected borrowed books routes
	borrowedBooks := e.Group("/borrowed-books", middleware.Auth)
	{
		borrowedBooks.POST("/borrow/:book_id", borrowedBooksHandler.BorrowBook)
		borrowedBooks.POST("/return/:id", borrowedBooksHandler.ReturnBook)
	}
}
