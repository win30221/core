package catch

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type customError struct {
	Code      string
	OutputMsg string
	LogMsg    string
	Stack     string
}

func getCallStack(min, max int) string {
	var callers = []string{}

	for i := min; i < max; i++ {
		_, path, line, ok := runtime.Caller(i)

		if !ok {
			break
		}

		var caller strings.Builder
		caller.WriteString(path)
		caller.WriteString(":")
		caller.WriteString(strconv.Itoa(line))

		callers = append(callers, caller.String())
	}

	return strings.Join(callers, " ")
}

// NewWitStack 用在自定義套件（如 core/storage/redis/redis.go）中的 error
func NewWitStack(code, outputMsg, logMsg string, startStack int) error {
	if startStack-1 < 0 {
		return customError{
			Code:      code,
			OutputMsg: outputMsg,
			LogMsg:    fmt.Sprintf("[warning]: start stack must greater than 0, %s", logMsg),
			Stack:     getCallStack(0, 10),
		}
	}

	// 找出呼叫 catch.New 的 function 名稱
	fName := ""
	pc, _, _, _ := runtime.Caller(startStack - 1)
	re := regexp.MustCompile(`.+\/.+\.(.+)`)
	match := re.FindStringSubmatch(runtime.FuncForPC(pc).Name())
	if len(match) > 1 {
		fName = match[1]
	}

	return customError{
		Code:      code,
		OutputMsg: outputMsg,
		LogMsg:    fmt.Sprintf("`%s`, %s", fName, logMsg),
		Stack:     getCallStack(startStack, startStack+4),
	}
}

// New 用在 http server 目前架構的層級中
func New(code, outputMsg, logMsg string) error {
	return NewWitStack(code, outputMsg, logMsg, 3)
}

func (e customError) Error() string {
	return fmt.Sprintf("Code: %s, OutputMsg: %s, LogMsg: %s, stack: %s", e.Code, e.OutputMsg, e.LogMsg, e.Stack)
}

// Info 給 http response package 用來輸出、列印訊息時使用
func (e customError) Info() (string, string, string, string) {
	return e.Code, e.OutputMsg, e.LogMsg, e.Stack
}

func CheckCustomError(err error) (e customError, ok bool) {
	e, ok = err.(customError)
	return
}

// ReplaceOutPutMsg 用在 gateway 的 repository 錯誤出現不同錯誤碼時會需要替換回傳的錯誤訊息
func ReplaceOutPutMsg(err error, msg string) error {
	e, ok := CheckCustomError(err)
	if !ok {
		return err
	}

	e.OutputMsg = msg

	return e
}
