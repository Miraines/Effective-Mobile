package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	deliv "Effective-Mobile/internal/delivery/http"
	"Effective-Mobile/internal/domain"
	"Effective-Mobile/internal/repository"
	"Effective-Mobile/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

type repoMini struct{ mock.Mock }

func (m *repoMini) Create(ctx context.Context, p *domain.Person) (int64, error) { return 0, nil }
func (m *repoMini) GetByID(ctx context.Context, id int64) (*domain.Person, error) {
	return &domain.Person{ID: id, Name: "Stub"}, nil
}
func (m *repoMini) List(ctx context.Context, f repository.ListFilter) ([]domain.Person, int, error) {
	return []domain.Person{{ID: 1, Name: "Stub"}}, 1, nil
}
func (m *repoMini) Update(ctx context.Context, p *domain.Person) error { return nil }
func (m *repoMini) Delete(ctx context.Context, id int64) error         { return nil }

type enrichMini struct{}

func (enrichMini) Enrich(context.Context, *domain.Person) error { return nil }

func newStubService(t *testing.T) *service.PeopleService {
	return service.NewPeopleService(&repoMini{}, enrichMini{}, zaptest.NewLogger(t).Sugar())
}

func TestHandler_List_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := deliv.NewHandler(newStubService(t))
	h.Register(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/people?limit=5", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "total")
}

func TestHandler_Get_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := deliv.NewHandler(newStubService(t))
	h.Register(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/people/42", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "\"id\":42")
}
