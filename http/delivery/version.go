package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/win30221/core/basic"
	"github.com/win30221/core/http/ctx"
	"github.com/win30221/core/http/response"
)

func Version(c *gin.Context) {
	ctx := ctx.New(c, c.Request.Context())

	response.OK(ctx, map[string]any{
		"Server":    basic.ServerName,
		"Host":      basic.Host,
		"Port":      basic.Port,
		"Version":   basic.Version,
		"BuildTime": basic.BuildTime,
		"Commit":    basic.Commit,
		"Consul":    basic.ConsulIP,
	})
}
