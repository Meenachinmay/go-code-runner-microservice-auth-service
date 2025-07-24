package company

import (
	"context"
	"go-code-runner-microservice/auth-service/internal/model"
)

type Repository interface {
	Create(ctx context.Context, c *model.Company) (*model.Company, error)
	GetByEmail(ctx context.Context, email string) (*model.Company, error)
	GetByID(ctx context.Context, id int) (*model.Company, error)
	GetCompanyByAPIKey(ctx context.Context, apiKey string) (*model.Company, error)
	UpdateAPIKey(ctx context.Context, id int, apiKey string) error
	UpdateClientID(ctx context.Context, id int, clientID string) error
}
