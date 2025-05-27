package service_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"Effective-Mobile/internal/domain"
	"Effective-Mobile/internal/service"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestEnrich_HTTPFailuresAreSoft(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://api.agify.io/?name=fail", httpmock.NewStringResponder(500, "boom"))
	httpmock.RegisterNoResponder(func(req *http.Request) (*http.Response, error) { return nil, context.DeadlineExceeded })

	e := service.NewEnricher(200 * time.Millisecond)

	p := &domain.Person{Name: "fail"}
	err := e.Enrich(context.Background(), p)

	assert.Error(t, err)
	assert.Nil(t, p.Age)
	assert.Nil(t, p.Gender)
	assert.Nil(t, p.CountryID)
}
