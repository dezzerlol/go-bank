package grpc

import (
	"context"
	"fmt"
	"go-bank/config"
	db "go-bank/db/sqlc"
	"go-bank/internal/token"
	"go-bank/pb"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

type GrpcServer struct {
	db         db.Store
	tokenMaker token.Maker
	cfg        config.Config

	pb.UnimplementedUserServiceServer
}

func NewGrpcServer(config config.Config, db db.Store) (*GrpcServer, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TOKEN_SYMMETRIC_KEY)

	if err != nil {
		return nil, fmt.Errorf("cant create server: %s", err)
	}

	server := &GrpcServer{
		db:         db,
		tokenMaker: tokenMaker,
		cfg:        config,
	}

	return server, nil
}

func (s *GrpcServer) Start() error {
	lis, err := net.Listen("tcp", s.cfg.GRPC_ADDR)

	if err != nil {
		return err
	}

	// Create grpc server
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, s)
	reflection.Register(grpcServer)

	// Start server
	log.Printf("Listening and serving GRPC on: %s \n", lis.Addr())
	err = grpcServer.Serve(lis)

	return err
}

func (s *GrpcServer) StartGateway() error {
	lis, err := net.Listen("tcp", s.cfg.SRV_ADDR)

	if err != nil {
		return err
	}

	// use snake_case as defined in proto file instead of camelCase
	protoNames := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	// Register grpc gateway
	grpcMux := runtime.NewServeMux(protoNames)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterUserServiceHandlerServer(ctx, grpcMux, s)

	if err != nil {
		return err
	}

	// Create http mux
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// Start server
	log.Printf("Listening and serving GRPC Gateway on: %s \n", s.cfg.SRV_ADDR)
	err = http.Serve(lis, mux)

	return err
}
