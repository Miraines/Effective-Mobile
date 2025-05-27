package service

import (
	"Effective-Mobile/internal/domain"
	"Effective-Mobile/internal/repository"
	"context"
	"errors"
	"go.uber.org/zap"
)

type PeopleService struct {
	repo     repository.PersonRepository
	enricher Enricher
	log      *zap.SugaredLogger
}

func NewPeopleService(r repository.PersonRepository, e Enricher, log *zap.SugaredLogger) *PeopleService {
	return &PeopleService{
		repo:     r,
		enricher: e,
		log:      log,
	}
}

func (s *PeopleService) Add(ctx context.Context, p *domain.Person) (int64, error) {
	if p.Name == "" {
		return 0, errors.New("name is required")
	}
	if err := s.enricher.Enrich(ctx, p); err != nil {
		s.log.Warnw("enrich failed", "err", err)
	}
	return s.repo.Create(ctx, p)
}

func (s *PeopleService) Get(ctx context.Context, id int64) (*domain.Person, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PeopleService) List(ctx context.Context, f repository.ListFilter) ([]domain.Person, int, error) {
	return s.repo.List(ctx, f)
}

func (s *PeopleService) Update(ctx context.Context, p *domain.Person) error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	if err := s.enricher.Enrich(ctx, p); err != nil {
		s.log.Warnw("enrich failed", "err", err)
	}
	return s.repo.Update(ctx, p)
}

func (s *PeopleService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
