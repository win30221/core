package validate

import "github.com/go-playground/validator/v10"

// request 的 struct filed 統一使用
var Validate = validator.New()

// delivery 在 c.ShouldBind() 之後使用 Validate.Struct() 來驗證傳入參數
func Struct(s interface{}) error {
	return Validate.Struct(s)
}
