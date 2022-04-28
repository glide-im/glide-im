package http_srv

import (
	"github.com/gin-gonic/gin"
	"go_im/im/api/comm"
	"net/http"
	"strings"
)

var authRouteGroup gin.IRoutes

const CtxKeyAuthInfo = "CTX_KEY_AUTH_INFO"

func useAuth() gin.IRoutes {
	if authRouteGroup == nil {
		authRouteGroup = g.Use(authMiddleware).Use(crosMiddleware())
	}
	return authRouteGroup
}

func authMiddleware(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	if authHeader == "" {
		context.Status(http.StatusUnauthorized)
		context.Abort()
		return
	}
	authHeader = strings.TrimPrefix(authHeader, "Bearer ")
	authInfo, err := comm.ParseJwt(authHeader)
	if err != nil {
		context.Status(http.StatusUnauthorized)
		context.Abort()
		return
	}
	context.Set(CtxKeyAuthInfo, authInfo)
	context.Next()
}
