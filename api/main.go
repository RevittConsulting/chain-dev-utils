package main

import (
	"context"
	"github.com/RevittConsulting/cdk-envs/config"
	"github.com/RevittConsulting/cdk-envs/internal/db"
	"github.com/RevittConsulting/cdk-envs/internal/jsonrpc"
	"github.com/RevittConsulting/cdk-envs/internal/server"
	"log"
)

func main() {
	server.Start(start)
}

func start(ctx context.Context, s *server.Server) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	database := db.New("boltdb.db")
	defer func() {
		if err := database.Close(); err != nil {
			log.Fatalf("Failed to close the database: %v", err)
		}
	}()

	rpcClient := jsonrpc.NewRPCClient(&cfg.RPC, database.Db)
	go rpcClient.GetMostRecentBlock()

	return s.Setup(ctx, cfg, database.Db)
}
