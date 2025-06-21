package app

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Ippolid/chat-server/internal/interceptor"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Ippolid/chat-server/internal/config"
	"github.com/Ippolid/chat-server/pkg/chatserver_v1"
	"github.com/Ippolid/platform_libary/pkg/closer"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// App представляет собой основное приложение
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

// NewApp создает новое приложение
func NewApp(ctx context.Context) (*App, error) {
	log.Println("Creating new App...")
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize dependencies: %w", err)
	}

	log.Println("App created successfully")
	return a, nil
}

// Run запускает приложение
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runGRPCServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load("./.env")
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	interceptors := grpc.ChainUnaryInterceptor(interceptor.ValidateInterceptor, a.serviceProvider.GetAuthInterceptor(ctx).AccessInterceptor)

	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()), interceptors)

	reflection.Register(a.grpcServer)

	chatserver_v1.RegisterChatV1Server(a.grpcServer, a.serviceProvider.NoteController(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	log.Printf("GRPC server is running on %s", a.serviceProvider.GRPCConfig().Address())

	list, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		log.Printf("failed to serve grpc server: %v", err)
		return err
	}

	return nil
}
