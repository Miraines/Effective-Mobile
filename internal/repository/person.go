package repository

import (
	"Effective-Mobile/internal/domain"
	"context"
)

type PersonRepository interface {
	Create(ctx context.Context, p *domain.Person) (int64, error)
	GetByID(ctx context.Context, id int64) (*domain.Person, error)
	List(ctx context.Context, filter ListFilter) ([]domain.Person, int, error)
	Update(ctx context.Context, p *domain.Person) error
	Delete(ctx context.Context, id int64) error
}

type ListFilter struct {
	Name    string
	Surname string
	AgeFrom *int
	AgeTo   *int
	Limit   int
	Offset  int
	SortBy  string
}
