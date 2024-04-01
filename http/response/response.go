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
	Data   interface{} `json:"data"`
	Status `json:"status"`
}

// Error 用在回傳值需要 data 的時候
func ErrorD(ctx ctx.Context, httpStatusCode int, data interface{}, err error) {
	var d interface{}

	if data != nil {
		d = data
	}

	customError, ok := catch.CheckCustomError(err)
	if !ok {
		ctx.GinContext.JSON(httpStatusCode, Response{
			Data: d,
			Status: Status{
				Code:      syserrno.Undefined,
				Message:   err.Error(),
				TraceCode: ctx.TraceCode,
				DateTime:  time.Now().In(basic.TimeZone).Format(time.RFC3339),
			},
		})

		ctx.GinContext.Error(errors.New(err.Error()))

		return
	}

	code, outputMsg, logMsg, stack := customError.Info()

	ctx.GinContext.JSON(httpStatusCode, Response{
		Data: d,
		Status: Status{
			Code:      code,
			TraceCode: ctx.TraceCode,
			Message:   outputMsg,
			DateTime:  time.Now().In(basic.TimeZone).Format(time.RFC3339),
		},
	})

	ctx.GinContext.Error(fmt.Errorf("%s, stack:%s", logMsg, stack))
}

// Error 用在回傳值沒有需要 data 的時候
func Error(ctx ctx.Context, httpStatusCode int, err error) {
	customError, ok := catch.CheckCustomError(err)
	if !ok {
		ctx.GinContext.JSON(httpStatusCode, Response{
			Status: Status{
				Code:      syserrno.Undefined,
				Message:   err.Error(),
				TraceCode: ctx.TraceCode,
				DateTime:  time.Now().In(basic.TimeZone).Format(time.RFC3339),
			},
		})

		ctx.GinContext.Error(errors.New(err.Error()))

		return
	}

	code, outputMsg, logMsg, stack := customError.Info()

	ctx.GinContext.JSON(httpStatusCode, Response{
		Status: Status{
			Code:      code,
			TraceCode: ctx.TraceCode,
			Message:   outputMsg,
			DateTime:  time.Now().In(basic.TimeZone).Format(time.RFC3339),
		},
	})

	ctx.GinContext.Error(fmt.Errorf("%s, stack:%s", logMsg, stack))
}

func OK(ctx ctx.Context, data interface{}) {
	res := &Response{
		Data: "Success",
		Status: Status{
			Code:      syserrno.OK,
			Message:   "Success",
			TraceCode: ctx.TraceCode,
			DateTime:  time.Now().In(basic.TimeZone).Format(time.RFC3339),
		},
	}

	if data != nil {
		res.Data = data
		ctx.GinContext.Set("result", data)
	}

	ctx.GinContext.JSON(http.StatusOK, res)
}

func BindParameterError(ctx ctx.Context, err error) {
	Error(ctx, http.StatusBadRequest, catch.New(
		syserrno.ValidParameter,
		fmt.Sprintf("bind parameter error: %v", err),
		err.Error(),
	))
}

func ValidParameterError(ctx ctx.Context, err error) {
	Error(ctx, http.StatusBadRequest, catch.New(
		syserrno.ValidParameter,
		fmt.Sprintf("validate parameter error: %v", err),
		err.Error(),
	))
}
