package main

import (
	"context"
	"net/http"
	"os"
	"time"

	pb "gc-buku/proto"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	e := echo.New()

	// Get gRPC server address from env
	grpcServer := os.Getenv("GRPC_SERVER")
	if grpcServer == "" {
		grpcServer = "server:50051" // default value
	}

	e.POST("/create_user", func(c echo.Context) error {
		type CreateUserRequest struct {
			Id       string `json:"id"`
			Username string `json:"username"`
			Password string `json:"password"`
		}
		req := new(CreateUserRequest)
		if err := c.Bind(req); err != nil {
			return c.String(http.StatusBadRequest, "invalid request: "+err.Error())
		}

		conn, err := grpc.Dial(grpcServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return c.String(http.StatusInternalServerError, "did not connect: "+err.Error())
		}
		defer conn.Close()

		client := pb.NewBookServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userResp, err := client.CreateUser(ctx, &pb.CreateUserRequest{
			User: &pb.User{
				Id:       req.Id,
				Username: req.Username,
				Password: req.Password,
			},
		})
		if err != nil {
			return c.String(http.StatusInternalServerError, "could not create user: "+err.Error())
		}

		return c.JSON(http.StatusOK, userResp.User)
	})

	e.GET("/get_user", func(c echo.Context) error {
		conn, err := grpc.Dial(grpcServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return c.String(http.StatusInternalServerError, "did not connect: "+err.Error())
		}
		defer conn.Close()

		client := pb.NewBookServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		getUserResp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: "1"})
		if err != nil {
			return c.String(http.StatusInternalServerError, "could not get user: "+err.Error())
		}

		return c.JSON(http.StatusOK, getUserResp.User)
	})

	e.Logger.Fatal(e.Start(":8081"))
}
