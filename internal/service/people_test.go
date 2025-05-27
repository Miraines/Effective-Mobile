package service_test

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"testing"

	"Effective-Mobile/internal/domain"
	"Effective-Mobile/internal/repository"
	"Effective-Mobile/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

type EnricherFunc func(context.Context, *domain.Person) error

func (f EnricherFunc) Enrich(ctx context.Context, p *domain.Person) error { return f(ctx, p) }

func zapLogger(t *testing.T) *zap.SugaredLogger { return zaptest.NewLogger(t).Sugar() }

type repoMock struct{ mock.Mock }

func (m *repoMock) Create(ctx context.Context, p *domain.Person) (int64, error) {
	args := m.Called(ctx, p)
	return int64(args.Int(0)), args.Error(1)
}
func (m *repoMock) GetByID(ctx context.Context, id int64) (*domain.Person, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Person), args.Error(1)
}
func (m *repoMock) List(ctx context.Context, f repository.ListFilter) ([]domain.Person, int, error) {
	args := m.Called(ctx, f)
	return args.Get(0).([]domain.Person), args.Int(1), args.Error(2)
}
func (m *repoMock) Update(ctx context.Context, p *domain.Person) error {
	return m.Called(ctx, p).Error(0)
}
func (m *repoMock) Delete(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}

func TestPeopleService_Add(t *testing.T) {
	r := new(repoMock)
	e := EnricherFunc(func(ctx context.Context, p *domain.Person) error { return nil })

	svc := service.NewPeopleService(r, e, zapLogger(t))

	p := &domain.Person{Name: "Alex"}
	r.On("Create", mock.Anything, p).Return(1, nil)

	id, err := svc.Add(context.Background(), p)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	r.AssertExpectations(t)
}

func TestPeopleService_Add_NoName(t *testing.T) {
	svc := service.NewPeopleService(new(repoMock), EnricherFunc(func(ctx context.Context, p *domain.Person) error { return nil }), zapLogger(t))
	_, err := svc.Add(context.Background(), &domain.Person{})
	assert.EqualError(t, err, "name is required")
}

func TestPeopleService_Update_EnrichFail(t *testing.T) {
	r := new(repoMock)
	enrichErr := errors.New("boom")

	svc := service.NewPeopleService(r, EnricherFunc(func(ctx context.Context, p *domain.Person) error { return enrichErr }), zapLogger(t))

	p := &domain.Person{ID: 1, Name: "Bob"}
	r.On("Update", mock.Anything, p).Return(nil)

	err := svc.Update(context.Background(), p)
	assert.NoError(t, err)
	r.AssertExpectations(t)
}
