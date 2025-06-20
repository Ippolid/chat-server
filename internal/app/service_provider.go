package app

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Ippolid/auth/pkg/auth_v1"
	"github.com/Ippolid/chat-server/internal/api/chatserver"
	"github.com/Ippolid/chat-server/internal/config"
	accessInterceptor "github.com/Ippolid/chat-server/internal/interceptor/access"
	"github.com/Ippolid/chat-server/internal/repository"
	"github.com/Ippolid/chat-server/internal/repository/auth"
	chat_server "github.com/Ippolid/chat-server/internal/repository/chat-server"
	"github.com/Ippolid/chat-server/internal/service"
	"github.com/Ippolid/chat-server/internal/service/access"
	chat_service "github.com/Ippolid/chat-server/internal/service/chatserver"
	"github.com/Ippolid/platform_libary/pkg/closer"
	"github.com/Ippolid/platform_libary/pkg/db"
	"github.com/Ippolid/platform_libary/pkg/db/pg"
	"github.com/Ippolid/platform_libary/pkg/db/transaction"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient         db.Client
	txManager        db.TxManager
	noteRepository   repository.ChatServerRepository
	accessRepository repository.Access
	authClient       auth_v1.AuthClient

	noteService   service.ChatServerService
	accessService service.Access

	noteController  *chatserver.Controller
	authInterceptor *accessInterceptor.AuthInterceptor
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

func (s *serviceProvider) GetAuthClient() auth_v1.AuthClient {
	if s.authClient == nil {
		// Вариант 1a: Загрузить клиентские сертификаты (если есть взаимная аутентификация)
		creds, err := credentials.NewClientTLSFromFile(
			"../../server_cert.pem", // путь к серверному сертификату
			"localhost",             // serverName (можно оставить пустым для localhost)
		)
		if err != nil {
			log.Fatalf("failed to load TLS credentials: %v", err)
		}

		authAddr := os.Getenv("AUTH_SERVER_ADDR")
		if authAddr == "" {
			authAddr = "host.docker.internal:50051"
		}

		conn, err := grpc.Dial(authAddr, grpc.WithTransportCredentials(creds))
		fmt.Println(conn, err)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		closer.Add(conn.Close)

		s.authClient = auth_v1.NewAuthClient(conn)
	}

	return s.authClient
}

func (s *serviceProvider) GetAccessRepository(_ context.Context) repository.Access {
	if s.accessRepository == nil {
		s.accessRepository = auth.NewAccessRepo(s.GetAuthClient())
	}

	return s.accessRepository
}

func (s *serviceProvider) GetAccessService(ctx context.Context) service.Access {
	if s.accessService == nil {
		s.accessService = access.NewAccessService(s.GetAccessRepository(ctx))
	}

	return s.accessService
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

func (s *serviceProvider) GetAuthInterceptor(ctx context.Context) accessInterceptor.AuthInterceptor {
	if s.authInterceptor == nil {
		s.authInterceptor = accessInterceptor.NewAuthInterceptor(s.GetAccessService(ctx))
	}

	return *s.authInterceptor
}
