package delivery

import (
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/win30221/core/basic"
	"github.com/win30221/core/http/middleware"
)

func SetBasicRouter(e *gin.Engine) (publicGroup, privateGroup *gin.RouterGroup) {
	publicGroup = e.Group("/" + basic.ServerName)
	privateGroup = e.Group("/"+basic.ServerName, middleware.ValidateToken(basic.SysToken))

	if basic.Site != "prd" {
		ginSwagger.WrapHandler(swaggerfiles.Handler,
			ginSwagger.URL(fmt.Sprintf("http://localhost:%s/%s/swagger/doc.json", basic.Port, basic.ServerName)),
			ginSwagger.DefaultModelsExpandDepth(-1),
		)

		// public group
		publicGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	// private group
	privateGroup.GET("/version", Version)
	return
}
