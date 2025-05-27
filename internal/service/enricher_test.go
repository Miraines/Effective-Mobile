package service_test

import (
	"context"
	"testing"
	"time"

	"Effective-Mobile/internal/domain"
	"Effective-Mobile/internal/service"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestEnrich_AllOk(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://api.agify.io/?name=sasha", httpmock.NewStringResponder(200, `{"age":42}`))
	httpmock.RegisterResponder("GET", "https://api.genderize.io/?name=sasha", httpmock.NewStringResponder(200, `{"gender":"male"}`))
	httpmock.RegisterResponder("GET", "https://api.nationalize.io/?name=sasha", httpmock.NewStringResponder(200, `{"country":[{"country_id":"FI","probability":0.9}]}`))

	e := service.NewEnricher(2 * time.Second)

	p := &domain.Person{Name: "sasha"}
	err := e.Enrich(context.Background(), p)
	assert.NoError(t, err)
	assert.Equal(t, 42, *p.Age)
	assert.Equal(t, "male", *p.Gender)
	assert.Equal(t, "FI", *p.CountryID)
}
