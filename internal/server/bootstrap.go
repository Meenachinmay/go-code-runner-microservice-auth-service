package server

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	companyauthpb "go-code-runner-microservice/auth-service/go-code-runner-microservice/auth/proto/company_auth/v1"
	"go-code-runner-microservice/auth-service/internal/config"
	"go-code-runner-microservice/auth-service/internal/grpc"
	"go-code-runner-microservice/auth-service/internal/platform/database"
	"go-code-runner-microservice/auth-service/internal/repository/company"
	companyService "go-code-runner-microservice/auth-service/internal/service/company"

	"github.com/joho/godotenv"
	grpcServer "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run() {
	logger := log.New(os.Stdout, "AUTH-SERVICE: ", log.LstdFlags|log.Lmicroseconds)
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("failed to load configuration: %v", err)
	}

	ctx := context.Background()
	dbpool, err := database.New(ctx, cfg.DBConnStr)
	if err != nil {
		logger.Fatalf("failed to connect to database: %v", err)
	}
	defer dbpool.Close()
	logger.Println("database connection pool established")

	logger.Println("checking for pending database migrations...")
	if err := database.Migrate(ctx, dbpool, "db/migrations", logger); err != nil {
		logger.Fatalf("migration failed: %v", err)
	}
	logger.Println("database is up-to-date")

	companyRepo := company.New(dbpool)
	companyService := companyService.New(companyRepo)

	grpcAddr := ":" + cfg.GrpcServerPort
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Fatalf("failed to listen on %s: %v", grpcAddr, err)
	}

	server := grpcServer.NewServer()

	// Register services
	companyAuthServer := grpc.NewCompanyAuthServer(companyService, logger)
	companyauthpb.RegisterCompanyAuthServiceServer(server, companyAuthServer)

	reflection.Register(server)

	go func() {
		logger.Printf("Starting gRPC server on %s", grpcAddr)
		if err := server.Serve(lis); err != nil {
			logger.Fatalf("failed to serve gRPC server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutdown signal received, initiating graceful shutdown...")

	server.GracefulStop()
	logger.Println("gRPC server stopped")
}
