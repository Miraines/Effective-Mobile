package service_test

import (
	"context"
	"testing"

	"Effective-Mobile/internal/domain"
	"Effective-Mobile/internal/repository"
	"Effective-Mobile/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

type repoStub struct{ mock.Mock }

func (m *repoStub) Create(ctx context.Context, p *domain.Person) (int64, error) {
	args := m.Called(ctx, p)
	return int64(args.Int(0)), args.Error(1)
}
func (m *repoStub) GetByID(ctx context.Context, id int64) (*domain.Person, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Person), args.Error(1)
}
func (m *repoStub) List(ctx context.Context, f repository.ListFilter) ([]domain.Person, int, error) {
	args := m.Called(ctx, f)
	return args.Get(0).([]domain.Person), args.Int(1), args.Error(2)
}
func (m *repoStub) Update(ctx context.Context, p *domain.Person) error {
	return m.Called(ctx, p).Error(0)
}
func (m *repoStub) Delete(ctx context.Context, id int64) error { return m.Called(ctx, id).Error(0) }

type enrichNoop struct{}

func (enrichNoop) Enrich(context.Context, *domain.Person) error { return nil }

func TestPeopleService_ListAndDelete(t *testing.T) {
	r := new(repoStub)
	svc := service.NewPeopleService(r, enrichNoop{}, zaptest.NewLogger(t).Sugar())

	expected := []domain.Person{{ID: 1, Name: "Amy"}}
	r.On("List", mock.Anything, mock.Anything).Return(expected, 1, nil)

	list, total, err := svc.List(context.Background(), repository.ListFilter{})
	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Equal(t, "Amy", list[0].Name)

	r.On("Delete", mock.Anything, int64(1)).Return(nil)
	assert.NoError(t, svc.Delete(context.Background(), 1))
	r.AssertExpectations(t)
}
