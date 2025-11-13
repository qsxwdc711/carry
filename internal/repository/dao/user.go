package dao

import (
	"context"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"sx-go/internal/domain"
	"sx-go/internal/web/middleware"
)

type UserDaoInterface interface {
	InsertOne(ctx context.Context, user User) (domain.User, error)
	FindByAccount(ctx context.Context, account string) (User, error)
}
type MongoUserDao struct {
	db *mongo.Collection
}

func NewUserDao(mongodb *mongo.Client) UserDaoInterface {
	database := viper.GetString("mongo.database")
	return &MongoUserDao{
		db: mongodb.Database(database).Collection("user"),
	}
}
func (dao *MongoUserDao) FindByAccount(ctx context.Context, account string) (User, error) {
	traceID := ""
	if v := ctx.Value(middleware.CtxTraceIDKey); v != nil {
		if ts, ok := v.(string); ok {
			traceID = ts
		}
	}
	zap.L().Info("dao.FindByUsername enter", zap.String("trace_id", traceID), zap.String("query", account))

	var user User
	err := dao.db.FindOne(ctx, bson.M{"account": account}).Decode(&user)
	if err != nil {

		zap.L().Info("dao.FindByUsername not found", zap.String("trace_id", traceID), zap.String("username", user.Name))
		return User{}, err
	}
	zap.L().Info("dao.FindByUsername found", zap.String("trace_id", traceID), zap.String("username", user.Account))
	return user, nil
}

func (dao *MongoUserDao) InsertOne(ctx context.Context, u User) (domain.User, error) {
	res, err := dao.db.InsertOne(ctx, &u)
	if err != nil {
		return domain.User{}, err
	}
	id := res.InsertedID.(primitive.ObjectID)
	return domain.User{Id: id}, nil
}

type User struct {
	Id       primitive.ObjectID `json:"id" bson:"_id,omitempty"` // id
	Account  string             `json:"account" bson:"account"`  // 账号
	Name     string             `json:"name" bson:"name"`        //姓名
	Phone    string             `json:"phone" bson:"phone"`
	Sex      string             `json:"sex" bson:"sex"`                       // 性别
	Password string             `json:"password" bson:"password"`             // 密码
	Avatar   string             `json:"avatar" bson:"avatar"`                 // 头像
	Role     primitive.ObjectID `json:"role,omitempty" bson:"role,omitempty"` // 角色外键
}
