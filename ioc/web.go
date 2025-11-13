package ioc

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"sx-go/internal/web"
	"sx-go/internal/web/middleware"
)

func InitGin(mdls []gin.HandlerFunc,
	userHdl *web.UserHandler,
) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterUserRoutes(server)
	//roleHdl.RegisterRoleRoutes(server)
	//centerHandler.RegisterRoutes(server)
	//groupHandler.RegisterRoutes(server)
	return server
}

func InitMiddlewares(db *mongo.Client) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.TraceMiddleware(),
		initCors(),
		middleware.NewLoginJWTMiddlewareBuilder().
			IgnorePaths("/users/register").
			IgnorePaths("/users/login").Build(),
		//ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
	}
}

func initCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
