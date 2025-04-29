package app

import (
	"context"
	"log"

	"github.com/Ippolid/chat-server/internal/api/chatserver"
	"github.com/Ippolid/chat-server/internal/client/db"
	"github.com/Ippolid/chat-server/internal/client/db/pg"
	"github.com/Ippolid/chat-server/internal/client/db/transaction"
	"github.com/Ippolid/chat-server/internal/closer"
	"github.com/Ippolid/chat-server/internal/config"
	"github.com/Ippolid/chat-server/internal/repository"
	chat_server "github.com/Ippolid/chat-server/internal/repository/chat-server"
	"github.com/Ippolid/chat-server/internal/service"
	chat_service "github.com/Ippolid/chat-server/internal/service/chatserver"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient       db.Client
	txManager      db.TxManager
	noteRepository repository.ChatServerRepository

	noteService service.ChatServerService

	noteController *chatserver.Controller
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}

		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) AuthRepository(ctx context.Context) repository.ChatServerRepository {
	if s.noteRepository == nil {
		s.noteRepository = chat_server.NewRepository(s.DBClient(ctx))
	}

	return s.noteRepository
}

func (s *serviceProvider) AuthService(ctx context.Context) service.ChatServerService {
	if s.noteService == nil {
		s.noteService = chat_service.NewService(
			s.AuthRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.noteService
}

func (s *serviceProvider) NoteController(ctx context.Context) *chatserver.Controller {
	if s.noteController == nil {
		s.noteController = chatserver.NewController(s.AuthService(ctx))
	}

	return s.noteController
}
