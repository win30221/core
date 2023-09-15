package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	"github.com/win30221/core/http/consts"
	"github.com/win30221/core/http/ctx"
	"github.com/win30221/core/http/response"
)

// 驗證內部服務間溝通用的 middleware
type Validator struct {
	validate *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

func ValidateToken(sysToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get(consts.HeaderSysToken) != sysToken {
			ctx := ctx.New(c, c.Request.Context())
			response.Error(ctx, http.StatusBadRequest, errors.New("validate system token error"))
			c.Abort()
		}
		c.Next()
	}
}
