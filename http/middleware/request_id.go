package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/win30221/core/http/consts"
)

func RequestIDMiddleware(c *gin.Context) {
	req := c.Request
	rid := req.Header.Get(consts.HeaderXRequestID)
	if rid == "" {
		rid = generateRequestID()
	}
	req.Header.Set(consts.HeaderXRequestID, rid)
	c.Next()
}
