package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	deliv "Effective-Mobile/internal/delivery/http"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestID_Middleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(deliv.RequestID())
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, c.GetString("request_id"))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	r.ServeHTTP(w, req)

	gotID := w.Body.String()
	assert.NotEmpty(t, gotID)
	assert.Equal(t, gotID, w.Header().Get("X-Request-ID"))
}
