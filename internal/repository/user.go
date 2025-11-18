package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"sx-go/internal/domain"
	"sx-go/internal/repository/dao"
	"sx-go/internal/web/middleware"
)

type UserRepoInterface interface {
	FindOneByAccount(ctx context.Context, account string) (domain.User, error)
	InsertOne(ctx context.Context, user domain.User) (domain.User, error)
	FindById(ctx context.Context, id primitive.ObjectID) (domain.User, error)
}
type UserRepository struct {
	dao dao.UserDaoInterface
}

func NewUserRepo(dao dao.UserDaoInterface) UserRepoInterface {
	return &UserRepository{
		dao: dao,
	}
}
func (repo *UserRepository) FindOneByAccount(ctx context.Context, account string) (domain.User, error) {
	user, err := repo.dao.FindByAccount(ctx, account)
	if err != nil {
		return domain.User{}, err
	}
	return toUserDomain(user), nil
}
func (repo *UserRepository) InsertOne(ctx context.Context, user domain.User) (domain.User, error) {
	return repo.dao.InsertOne(ctx, dao.User{
		Account:  user.Account,
		Password: user.Password,
		Name:     user.Name,
		Phone:    user.Phone,
		Sex:      user.Sex,
	})
}

func (repo *UserRepository) FindById(ctx context.Context, id primitive.ObjectID) (domain.User, error) {
	traceVal := ctx.Value(middleware.CtxTraceIDKey)
	trace := ""
	if s, ok := traceVal.(string); ok {
		trace = s
	}
	// 然后在日志中使用
	zap.L().Info("repo.FindById called", zap.String("trace_id", trace), zap.String("id", id.Hex()))
	user, err := repo.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return toUserDomain(user), nil
}
func toUserDomain(user dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Account:  user.Account,
		Name:     user.Name,
		Phone:    user.Phone,
		Sex:      user.Sex,
		Password: user.Password,
		//Avatar:   user.Avatar,
		//Role:     user.Role,
	}
}
