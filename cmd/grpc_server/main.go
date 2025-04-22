package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Ippolid/chat-server/internal/config"
	"github.com/Ippolid/chat-server/internal/server"
	"github.com/Ippolid/chat-server/pkg/chatserver_v1"
	"github.com/jackc/pgx/v5"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50052

func main() {
	ctx := context.Background()

	err := config.Load("./.env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dsn, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	con, err := pgx.Connect(ctx, dsn.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func(con *pgx.Conn, ctx context.Context) {
		err := con.Close(ctx)
		if err != nil {
			log.Printf("failed to close connection: %v", err)
		}
	}(con, ctx)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	serv := server.NewServer(con)
	s := grpc.NewServer()
	reflection.Register(s)
	chatserver_v1.RegisterChatV1Server(s, serv)

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
