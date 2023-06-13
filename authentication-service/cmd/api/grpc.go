package main

import (
	"authentication-service/data"
	auth "authentication-service/proto"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	Models data.Models
}

func (authServer *AuthServer) Authenticate(ctx context.Context, req *auth.AuthRequest) (*auth.AuthResponse, error) {
	input := req.GetAuthEntry()
	user := data.User{
		Email:    input.Email,
		Password: input.Password,
	}
	exist, err := authServer.Models.User.GetByEmail(user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = authServer.Models.User.Insert(user)
			if err != nil {
				res := &auth.AuthResponse{Result: "failed"}
				return res, err
			}
			res := &auth.AuthResponse{Result: "authenticated!"}
			return res, nil
		} else {
			res := &auth.AuthResponse{Result: "failed"}
			return res, err
		}
	}
	if bcrypt.CompareHashAndPassword([]byte(exist.Password), []byte(user.Password)) != nil {
		res := &auth.AuthResponse{Result: "password does not match!"}
		return res, errors.New("password error")
	}
	// return response
	res := &auth.AuthResponse{Result: "authenticated!"}
	return res, nil
}
func (app *Config) grpcListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
	s := grpc.NewServer()
	auth.RegisterAuthServiceServer(s, &AuthServer{Models: app.Models})
	log.Printf("gRPC Server started on port %s", grpcPort)
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
