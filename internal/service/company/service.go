package company

import (
	"context"
	"go-code-runner-microservice/auth-service/internal/model"
)

type Service interface {
	Register(ctx context.Context, name, email, password string) (*model.Company, error)
	Login(ctx context.Context, email, password string) (*model.Company, string, error)
	GenerateAPIKey(ctx context.Context, companyID int) (string, error)
	GenerateClientID(ctx context.Context, companyID int) (string, error)
}