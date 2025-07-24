package company

import (
	"context"
	"errors"
	"go-code-runner-microservice/auth-service/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) Repository {
	return &repo{db: db}
}

func (r *repo) GetCompanyByAPIKey(ctx context.Context, apiKey string) (*model.Company, error) {
	query := `
        SELECT id, name, email, password_hash, api_key, client_id, created_at, updated_at
        FROM companies
        WHERE api_key = $1`

	var company model.Company
	var apiKeyPtr, clientIDPtr *string
	err := r.db.QueryRow(ctx, query, apiKey).Scan(
		&company.ID,
		&company.Name,
		&company.Email,
		&company.PasswordHash,
		&apiKeyPtr,
		&clientIDPtr,
		&company.CreatedAt,
		&company.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	company.APIKey = apiKeyPtr
	company.ClientID = clientIDPtr

	return &company, nil
}

func (r *repo) Create(ctx context.Context, c *model.Company) (*model.Company, error) {
	query := `INSERT INTO companies (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, c.Name, c.Email, c.PasswordHash).
		Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *repo) GetByEmail(ctx context.Context, email string) (*model.Company, error) {
	var c model.Company
	var apiKeyPtr, clientIDPtr *string
	query := `SELECT id, name, email, password_hash, api_key, client_id, created_at, updated_at FROM companies WHERE email = $1`
	err := r.db.QueryRow(ctx, query, email).
		Scan(&c.ID, &c.Name, &c.Email, &c.PasswordHash, &apiKeyPtr, &clientIDPtr, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	
	c.APIKey = apiKeyPtr
	c.ClientID = clientIDPtr
	
	return &c, nil
}

func (r *repo) GetByID(ctx context.Context, id int) (*model.Company, error) {
	var c model.Company
	var apiKeyPtr, clientIDPtr *string
	query := `SELECT id, name, email, password_hash, api_key, client_id, created_at, updated_at FROM companies WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).
		Scan(&c.ID, &c.Name, &c.Email, &c.PasswordHash, &apiKeyPtr, &clientIDPtr, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	
	c.APIKey = apiKeyPtr
	c.ClientID = clientIDPtr
	
	return &c, nil
}

func (r *repo) UpdateAPIKey(ctx context.Context, id int, apiKey string) error {
	tag, err := r.db.Exec(ctx, `UPDATE companies SET api_key = $1, updated_at = now() WHERE id = $2`, apiKey, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("company not found")
	}
	return nil
}

func (r *repo) UpdateClientID(ctx context.Context, id int, clientID string) error {
	tag, err := r.db.Exec(ctx, `UPDATE companies SET client_id = $1, updated_at = now() WHERE id = $2`, clientID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("company not found")
	}
	return nil
}