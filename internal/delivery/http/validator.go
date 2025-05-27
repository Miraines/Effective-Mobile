package http

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func badValidation(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out := make(map[string]string, len(errs))
	for _, fe := range errs {
		out[fe.Field()] = fe.Tag()
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error":   "validation",
		"details": out,
	})
}

func init() {
	_ = validate.RegisterValidation("alpharus", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^[\p{L}]+$`).MatchString(fl.Field().String())
	})
}
