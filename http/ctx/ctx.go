package ctx

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/win30221/core/http/consts"
)

// 如果是在 rmq 或 cron-job 中使用，就不會有 gin 的 context，因此直接使用 Context 的 struct 來建立 ctx
type Context struct {
	// 在 response 的時候 callback
	GinContext *gin.Context
	// request
	Context   context.Context
	TraceCode string
}

func New(c *gin.Context, ctx context.Context) *Context {
	return &Context{
		GinContext: c,
		Context:    ctx,
		TraceCode:  c.Request.Header.Get(consts.HeaderXRequestID),
	}
}

func NewEmpty() *Context {
	return &Context{
		GinContext: &gin.Context{},
		Context:    context.Background(),
	}
}
