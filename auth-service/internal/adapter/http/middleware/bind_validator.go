package middleware

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/thanhnamdk2710/auth-service/internal/shared/response"
	"github.com/thanhnamdk2710/auth-service/internal/shared/validation"
)

func BindAndValidate[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T

		if err := c.ShouldBindJSON(&req); err != nil {
			if verrs, ok := err.(validator.ValidationErrors); ok {
				errors := validation.ConvertValidationErrors(verrs)
				response.Error(c, http.StatusBadRequest, "INPUT_INVALID", errors)
				c.Abort()
				return
			}
			response.Error(c, http.StatusBadRequest, "INVALID_JSON", "Invalid request body format")
			c.Abort()
			return
		}

		if v := reflect.ValueOf(req); v.Kind() == reflect.Struct {
			method := v.MethodByName("Validate")
			if method.IsValid() {
				results := method.Call(nil)
				if len(results) > 0 && !results[0].IsNil() {
					err := results[0].Interface().(error)
					response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
					c.Abort()
					return
				}
			}
		}

		c.Set("requestBody", req)

		c.Next()
	}
}
