package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlphaRusRule(t *testing.T) {
	ok := "Сашка"
	bad := "Sasha1"

	assert.True(t, validate.Var(ok, "alpharus") == nil)
	assert.Error(t, validate.Var(bad, "alpharus"))
}
