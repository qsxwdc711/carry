package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"sx-go/internal/domain"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}
func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	//添加忽略路径
	l.paths = append(l.paths, path)
	return l
}
func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		//忽略路径
		//如果是登录和注册不需要校验
		for _, path := range l.paths {
			if c.Request.URL.Path == path {
				return
			}
		}
		//我现在用jwt来检验
		tokenHeader := c.GetHeader("Authorization")
		if tokenHeader == "" {
			tokenHeader = c.Request.URL.Query().Get("Authorization")
		}
		if tokenHeader == "" {
			//说明没有登录
			c.AbortWithStatus(401)
			return
		}
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			c.AbortWithStatus(401)
			return
		}
		//获取token
		var claims = domain.UserClaims{}
		//获取除了Bearer那一坨字符串
		tokenStr := segs[1]
		secret := viper.GetString("general.jwt")
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			//密钥
			return []byte(secret), nil
		})
		//有错误 过期了 uid为初始值
		if err != nil || !token.Valid || claims.Uid == primitive.NilObjectID {
			//解析失败 没有登录
			c.AbortWithStatus(401)
			return
		}
		//如果没有过期1
		//每1分钟刷新一次
		//now := time.Now()
		//if claims.ExpiresAt.Sub(now) < time.Hour*23+time.Minute*59 {
		//	//claims是用来生成token的有关对象,设置一个新的token 7天有效期
		//	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24))
		//	tokenStr, err := token.SignedString([]byte(secret))
		//	if err != nil {
		//		log.Println("续约失败")
		//		return
		//	}
		//	//重新设置一个token
		//	c.Header("x-jwt-token", tokenStr)
		//}
		//err为空即解析成功用户可以登录,每次请求接口时都将Authorization解析后并设置claims,里面有id等
		c.Set("claims", claims)

	}
}
