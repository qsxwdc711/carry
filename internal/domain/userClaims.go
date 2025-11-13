package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserClaims struct {
	jwt.RegisteredClaims
	//声明你自己要放在token里面的数据
	Uid primitive.ObjectID
}
