package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"sx-go/internal/domain"
	"sx-go/internal/service"
	"sx-go/internal/web/middleware"
)

type RegisterReq struct {
	Account  string `json:"account" bson:"account"` // 账号
	Name     string `json:"name" bson:"name"`       //姓名
	Phone    string `json:"phone" bson:"phone"`
	Sex      string `json:"sex" bson:"sex"`           // 性别
	Password string `json:"password" bson:"password"` // 密码
	//Avatar   string `json:"avatar" bson:"avatar"`     // 头像
	//Role     string `json:"role" bson:"role"`
}
type LoginReq struct {
	Account  string `json:"account" bson:"account"`
	Password string `json:"password" bson:"password"`
}
type UserHandler struct {
	svc service.UserServiceInterface
}

func NewUserHandler(svc service.UserServiceInterface) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (u *UserHandler) RegisterUserRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/register", u.Register)
	ug.POST("/login", u.login)
	ug.GET("/profile", u.profile)
}

func (u *UserHandler) Register(ctx *gin.Context) {
	traceID, _ := ctx.Get(middleware.CtxTraceIDKey)
	traceStr := ""
	if s, ok := traceID.(string); ok {
		traceStr = s
	}

	zap.L().Info("handler.Register enter", zap.String("trace_id", traceStr))
	var req RegisterReq
	if err := ctx.ShouldBind(&req); err != nil {
		zap.L().Warn("handler.Register bind error", zap.Error(err), zap.String("trace_id", traceStr))
		ctx.JSON(http.StatusBadRequest, "参数格式错误")
		return
	}

	data, err := u.svc.Register(ctx, domain.User{
		Account:  req.Account,
		Name:     req.Name,
		Phone:    req.Phone,
		Sex:      req.Sex,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrRegisted) {
		ctx.JSON(http.StatusConflict, "用户已经注册了")
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	zap.L().Info("handler.Register success", zap.String("user", req.Account), zap.String("trace_id", traceStr))
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    data,
		"message": "注册成功",
	})
	return
}

func (u *UserHandler) login(ctx *gin.Context) {
	traceID, _ := ctx.Get(middleware.CtxTraceIDKey)
	traceStr := ""
	if s, ok := traceID.(string); ok {
		traceStr = s
	}

	zap.L().Info("handler.Login enter", zap.String("trace_id", traceStr))
	// 1. 初始化请求参数结构体
	var req LoginReq

	// 2. 从请求中解析参数（以 JSON 格式为例）
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 解析失败（如参数缺失、格式错误），返回 400 错误
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误：" + err.Error(),
		})
		return
	}
	data, err := u.svc.Login(ctx, req.Account, req.Password)

	if errors.Is(err, service.ErrNotRegisted) {
		ctx.JSON(http.StatusUnauthorized, "用户未注册")
		return
	}
	if errors.Is(err, service.ErrFalsePassword) {
		ctx.JSON(http.StatusForbidden, "密码错误")
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "系统错误")
		return
	}
	zap.L().Info("handler.Login success", zap.String("user", req.Account), zap.String("trace_id", traceStr))
	ctx.JSON(http.StatusOK, gin.H{
		"data":    data,
		"code":    200,
		"message": "登录成功",
	})
}

func (u *UserHandler) profile(ctx *gin.Context) {
	claimsVal, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "为获取到用户信息",
		})
		return
	}
	fmt.Println(claimsVal)
	claims, ok := claimsVal.(domain.UserClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "用户信息格式错误",
		})
		return
	}
	uid := claims.Uid

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"user_id": uid,
	})

}
