package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"sx-go/internal/domain"
	"sx-go/internal/repository"
	"sx-go/internal/web/middleware"
	"time"
)

var (
	ErrRegisted      = errors.New("用户已经注册了")
	ErrNotRegisted   = errors.New("用户没有注册")
	ErrFalsePassword = errors.New("密码错误")
)

type UserServiceInterface interface {
	Register(ctx context.Context, u domain.User) (any, error)
	Login(ctx context.Context, account string, password string) (any, error)
}
type UserService struct {
	repo repository.UserRepoInterface
}

func NewUserService(repo repository.UserRepoInterface) UserServiceInterface {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) Register(ctx context.Context, u domain.User) (any, error) {
	traceVal := ""
	if v := ctx.Value(middleware.CtxTraceIDKey); v != nil {
		if ts, ok := v.(string); ok {
			traceVal = ts
		}
	}
	zap.L().Info("service.Register enter", zap.String("trace_id", traceVal), zap.String("username", u.Account))
	_, err := svc.repo.FindOneByAccount(ctx, u.Account)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {

		zap.L().Error("service.Register repo.FindByUsername error", zap.Error(err), zap.String("trace_id", traceVal))
		return nil, err
	}
	if err == nil {
		zap.L().Warn("service.Register user exists", zap.String("username", u.Account), zap.String("trace_id", traceVal))
		return nil, ErrRegisted
	}
	//密码加密
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {

		return nil, err
	}
	u.Password = string(hash)
	//入库
	res, err := svc.repo.InsertOne(ctx, u)
	if err != nil {
		zap.L().Error("service.Register repo.Create error", zap.Error(err), zap.String("trace_id", traceVal))
		return nil, err
	}
	//生成token
	token, err := createToken(res.Id)
	res.Token = token
	zap.L().Info("service.Register success", zap.String("username", u.Account), zap.String("trace_id", traceVal))
	return res, nil
}
func (svc *UserService) Login(ctx context.Context, account string, password string) (any, error) {
	traceVal := ""
	if v := ctx.Value(middleware.CtxTraceIDKey); v != nil {
		if ts, ok := v.(string); ok {
			traceVal = ts
		}
	}
	zap.L().Info("service.Login enter", zap.String("trace_id", traceVal), zap.String("username", account))
	user, err := svc.repo.FindOneByAccount(ctx, account)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotRegisted
	}
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrFalsePassword
	}
	//生成token
	tokenString, err := createToken(user.Id)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	zap.L().Info("service.Login success", zap.String("username", user.Account), zap.String("trace_id", traceVal))
	return gin.H{
		"token": "bearer" + " " + tokenString,
		"user":  user,
	}, nil
}

//	func (svc *UserService) profile(c *gin.Context) {
//		uid, err := svc.repo.FindById(ctx, id)
//		if err != nil {
//
//		}
//		username, err := svc.repo.FindOneByAccount()
//	}
func createToken(Id primitive.ObjectID) (tokenString string, err error) {
	// 创建一个我们自己的声明
	secret := viper.GetString("general.jwt")
	claims := domain.UserClaims{
		//设置参数
		RegisteredClaims: jwt.RegisteredClaims{
			//设置7天的过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Uid: Id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//加密
	tokenStr, err := token.SignedString([]byte(secret))
	return tokenStr, err
}
