package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID `json:"id,omitempty"` // id
	Account  string             `json:"account"`      // 账号
	Name     string             `json:"name"`         //姓名
	Phone    string             `json:"phone"`
	Sex      string             `json:"sex" `     // 性别
	Password string             `json:"password"` // 密码
	Token    string             `json:"token"`
}
