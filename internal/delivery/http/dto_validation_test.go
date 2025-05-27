package http

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestCreatePersonRequest_Validation(t *testing.T) {
	valid := CreatePersonRequest{
		Name:    "Иван",
		Surname: "Петров",
	}
	assert.NoError(t, validate.Struct(valid))

	invalid := CreatePersonRequest{
		Name:    "Ivan123",
		Surname: "",
	}
	err := validate.Struct(invalid)
	assert.Error(t, err)

	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationErrors, got %T: %v", err, err)
	}
	tags := make(map[string]string, len(ve))
	for _, fe := range ve {
		tags[fe.Field()] = fe.Tag()
	}

	assert.Equal(t, "alpharus", tags["Name"])
	assert.Equal(t, "required", tags["Surname"])
}
