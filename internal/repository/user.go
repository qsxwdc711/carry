package repository

import (
	"context"
	"sx-go/internal/domain"
	"sx-go/internal/repository/dao"
)

type UserRepoInterface interface {
	FindOneByAccount(ctx context.Context, account string) (domain.User, error)
	InsertOne(ctx context.Context, user domain.User) (domain.User, error)
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
