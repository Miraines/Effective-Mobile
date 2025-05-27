package http_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	deliv "Effective-Mobile/internal/delivery/http"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Create_Validate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := deliv.NewHandler(nil)
	h.Register(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/people", bytes.NewBufferString(`{"name":"","surname":""}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "validation")
}
