package response

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/win30221/core/basic"
	"github.com/win30221/core/http/catch"
	"github.com/win30221/core/http/ctx"
	"github.com/win30221/core/syserrno"
)

// Status - 回傳 client 的狀態欄位
type Status struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	DateTime  string `json:"dateTime"`
	TraceCode string `json:"traceCode"`
}

type Response struct {
	Data   any `json:"data"`
	Status `json:"status"`
}

// Error 用在回傳值需要 data 的時候
func ErrorD(c *ctx.Context, httpStatusCode int, data any, err error) {
	var d any

	if data != nil {
		d = data
	}

	customError, ok := catch.CheckCustomError(err)
	if !ok {
		c.GinContext.JSON(httpStatusCode, Response{
			Data: d,
			Status: Status{
				Code:      syserrno.Undefined,
				Message:   err.Error(),
				TraceCode: c.TraceCode,
				DateTime:  time.Now().In(basic.TimeZone).Format(time.RFC3339),
			},
		})

		c.GinContext.Error(errors.New(err.Error()))

		return
	}

	code, outputMsg, logMsg, stack := customError.Info()

	c.GinContext.JSON(httpStatusCode, Response{
		Data: d,
		Status: Status{
			Code:      code,
			TraceCode: c.TraceCode,
			Message:   outputMsg,
			DateTime:  time.Now().In(basic.TimeZone).Format(time.RFC3339),
		},
	})

	c.GinContext.Error(fmt.Errorf("%s, stack:%s", logMsg, stack))
}

// Error 用在回傳值沒有需要 data 的時候
func Error(c *ctx.Context, httpStatusCode int, err error) {
	customError, ok := catch.CheckCustomError(err)
	if !ok {
		c.GinContext.JSON(httpStatusCode, Response{
			Status: Status{
				Code:      syserrno.Undefined,
				Message:   err.Error(),
				TraceCode: c.TraceCode,
				DateTime:  time.Now().In(basic.TimeZone).Format(time.RFC3339),
			},
		})

		c.GinContext.Error(errors.New(err.Error()))

		return
	}

	code, outputMsg, logMsg, stack := customError.Info()

	c.GinContext.JSON(httpStatusCode, Response{
		Status: Status{
			Code:      code,
			TraceCode: c.TraceCode,
			Message:   outputMsg,
			DateTime:  time.Now().In(basic.TimeZone).Format(time.RFC3339),
		},
	})

	c.GinContext.Error(fmt.Errorf("%s, stack:%s", logMsg, stack))
}

func OK(c *ctx.Context, data any) {
	res := &Response{
		Data: "Success",
		Status: Status{
			Code:      syserrno.OK,
			Message:   "Success",
			TraceCode: c.TraceCode,
			DateTime:  time.Now().In(basic.TimeZone).Format(time.RFC3339),
		},
	}

	if data != nil {
		res.Data = data
		c.GinContext.Set("result", data)
	}

	c.GinContext.JSON(http.StatusOK, res)
}

func BindParameterError(c *ctx.Context, err error) {
	Error(c, http.StatusBadRequest, catch.New(
		syserrno.ValidParameter,
		fmt.Sprintf("bind parameter error: %v", err),
		err.Error(),
	))
}

func ValidParameterError(c *ctx.Context, err error) {
	Error(c, http.StatusBadRequest, catch.New(
		syserrno.ValidParameter,
		fmt.Sprintf("validate parameter error: %v", err),
		err.Error(),
	))
}
