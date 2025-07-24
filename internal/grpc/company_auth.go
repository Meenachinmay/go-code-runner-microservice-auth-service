package grpc

import (
	"context"
	"log"

	companyauthpb "go-code-runner-microservice/auth-service/go-code-runner-microservice/auth/proto/company_auth/v1"
	"go-code-runner-microservice/auth-service/internal/model"
	"go-code-runner-microservice/auth-service/internal/service/company"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CompanyAuthServer struct {
	companyauthpb.UnimplementedCompanyAuthServiceServer
	service company.Service
	logger  *log.Logger
}

func NewCompanyAuthServer(service company.Service, logger *log.Logger) *CompanyAuthServer {
	return &CompanyAuthServer{
		service: service,
		logger:  logger,
	}
}

func (s *CompanyAuthServer) Register(ctx context.Context, req *companyauthpb.RegisterRequest) (*companyauthpb.RegisterResponse, error) {
	comp, err := s.service.Register(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		s.logger.Printf("Failed to register company: %v", err)
		errMsg := err.Error()
		return &companyauthpb.RegisterResponse{
			Success: false,
			Error:   &errMsg,
		}, nil
	}

	return &companyauthpb.RegisterResponse{
		Success: true,
		Company: convertCompanyToProto(comp),
	}, nil
}

func (s *CompanyAuthServer) Login(ctx context.Context, req *companyauthpb.LoginRequest) (*companyauthpb.LoginResponse, error) {
	comp, token, err := s.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		s.logger.Printf("Failed to login: %v", err)
		errMsg := err.Error()
		return &companyauthpb.LoginResponse{
			Success: false,
			Error:   &errMsg,
		}, nil
	}

	return &companyauthpb.LoginResponse{
		Success: true,
		Company: convertCompanyToProto(comp),
		Token:   &token,
	}, nil
}

func (s *CompanyAuthServer) GenerateAPIKey(ctx context.Context, req *companyauthpb.GenerateAPIKeyRequest) (*companyauthpb.GenerateAPIKeyResponse, error) {
	apiKey, err := s.service.GenerateAPIKey(ctx, int(req.CompanyId))
	if err != nil {
		s.logger.Printf("Failed to generate API key: %v", err)
		errMsg := err.Error()
		return &companyauthpb.GenerateAPIKeyResponse{
			Success: false,
			Error:   &errMsg,
		}, nil
	}

	return &companyauthpb.GenerateAPIKeyResponse{
		Success: true,
		ApiKey:  &apiKey,
	}, nil
}

func (s *CompanyAuthServer) GenerateClientID(ctx context.Context, req *companyauthpb.GenerateClientIDRequest) (*companyauthpb.GenerateClientIDResponse, error) {
	clientID, err := s.service.GenerateClientID(ctx, int(req.CompanyId))
	if err != nil {
		s.logger.Printf("Failed to generate client ID: %v", err)
		errMsg := err.Error()
		return &companyauthpb.GenerateClientIDResponse{
			Success: false,
			Error:   &errMsg,
		}, nil
	}

	return &companyauthpb.GenerateClientIDResponse{
		Success:  true,
		ClientId: &clientID,
	}, nil
}

func convertCompanyToProto(comp *model.Company) *companyauthpb.Company {
	protoComp := &companyauthpb.Company{
		Id:        int32(comp.ID),
		Name:      comp.Name,
		Email:     comp.Email,
		CreatedAt: timestamppb.New(comp.CreatedAt),
		UpdatedAt: timestamppb.New(comp.UpdatedAt),
	}

	if comp.APIKey != nil {
		apiKey := *comp.APIKey
		protoComp.ApiKey = &apiKey
	}

	if comp.ClientID != nil {
		clientID := *comp.ClientID
		protoComp.ClientId = &clientID
	}

	return protoComp
}
