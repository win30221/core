package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/win30221/core/http/consts"
	"github.com/win30221/core/utils"
)

func RequestIdMiddleware(c *gin.Context) {
	req := c.Request
	rid := req.Header.Get(consts.HeaderXRequestId)
	if rid == "" {
		rid = utils.GenerateRequestId()
	}
	req.Header.Set(consts.HeaderXRequestId, rid)
	c.Next()
}
