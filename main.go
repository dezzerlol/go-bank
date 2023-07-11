package main

import (
	"database/sql"
	"go-bank/cmd/api"
	"go-bank/cmd/grpc"
	"go-bank/config"
	"log"

	db "go-bank/db/sqlc"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load(".")
	if err != nil {
		log.Fatal("Cant read config file: ", err)
	}

	conn, err := sql.Open("postgres", cfg.DB_DSN)
	if err != nil {
		log.Fatal("Cant connect to database: ", err)
	}

	db := db.NewStore(conn)

	go func() {
		startGrpcServer(cfg, db)
	}()

	startGatewayServer(cfg, db)

	//startHTTPServer(cfg, db)
}

func startHTTPServer(cfg config.Config, db db.Store) {
	httpServer, err := api.NewHTTPServer(cfg, db)
	if err != nil {
		log.Fatal("Cant create http server: ", err)
	}

	err = httpServer.Start(cfg.SRV_ADDR)
	if err != nil {
		log.Fatal("Cant start http server: ", err)
	}
}

func startGrpcServer(cfg config.Config, db db.Store) {
	grpcServer, err := grpc.NewGrpcServer(cfg, db)
	if err != nil {
		log.Fatal("Cant create grpc server: ", err)
	}

	err = grpcServer.Start()
	if err != nil {
		log.Fatal("Cant start grpc server: ", err)
	}
}

func startGatewayServer(cfg config.Config, db db.Store) {
	grpcServer, err := grpc.NewGrpcServer(cfg, db)
	if err != nil {
		log.Fatal("Cant create grpc gateway server: ", err)
	}

	err = grpcServer.StartGateway()
	if err != nil {
		log.Fatal("Cant start grpc gateway server: ", err)
	}
}
