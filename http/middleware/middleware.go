package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/win30221/core/basic"
	"github.com/win30221/core/http/consts"
	"go.uber.org/zap"
)

const (
	SQLLogs = "sqlLogs"
)

func Log() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		RequestIDMiddleware,
		ginLogger(),
	}
}

func ginLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(SQLLogs, []string{})
		reckon := time.Now()
		c.Next()

		if excludePath(c.Request.RequestURI) {
			return
		}

		lastErr := c.Errors.Last()
		var err error
		if lastErr != nil {
			err = lastErr.Err
		}

		fs := []zap.Field{}
		fs = append(fs, basicFields(c, reckon)...)
		fs = append(fs, dumpHeader(c.Request)...)
		fs = append(fs, dumpForm(c.Request)...)
		if err != nil {
			zap.L().Error(err.Error(), fs...)
			return
		}

		if time.Since(reckon).Milliseconds() >= int64(basic.RequestLatencyThrottle) {
			zap.L().Warn(fmt.Sprintf("over latency %d(ms)", basic.RequestLatencyThrottle), fs...)
			return
		}

		zap.L().Info("", fs...)
	}
}

// excludePath 過濾不必要輸出的路徑
func excludePath(path string) bool {
	return strings.HasSuffix(path, "/version") ||
		strings.HasSuffix(path, "/ping")
}

// basicFields 記錄一些必要的資訊
func basicFields(c *gin.Context, reckon time.Time) (res []zap.Field) {
	res = []zap.Field{
		zap.String("traceCode", c.Request.Header.Get(consts.HeaderXRequestID)),
		zap.String("method", c.Request.Method),
		zap.String("uri", c.Request.RequestURI),
		zap.Int("status", c.Writer.Status()),
		zap.String("remoteHost", c.Request.RemoteAddr),
		zap.Duration("latency", time.Since(reckon)),
	}

	if basic.PrintDetail {
		sqlLogs, _ := c.Get(SQLLogs)
		result, _ := c.Get("result")
		res = append(res,
			zap.Any("sqlLogs", sqlLogs),
			zap.Any("result", result),
		)
	}

	return
}

// dumpHeader 顯示標頭
func dumpHeader(req *http.Request) (fs []zap.Field) {
	fs = []zap.Field{zap.Skip()}
	if req == nil {
		return
	}

	header := make(map[string][]string)
	v, ok := req.Header[consts.HeaderSysToken]
	if ok {
		header[consts.HeaderSysToken] = v
	}

	v, ok = req.Header[consts.HeaderAuthorization]
	if ok {
		header[consts.HeaderAuthorization] = v
	}

	fs = append(fs, zap.Any("HEADER", header))

	return
}

// dumpForm 顯示 request 參數
func dumpForm(req *http.Request) (fs []zap.Field) {
	fs = []zap.Field{zap.Skip()}
	if req == nil {
		return
	}

	err := req.ParseForm()
	if err != nil {
		zap.L().Error(err.Error())
		return
	}
	fs = append(fs, zap.Any("FORM", req.Form))

	return
}
