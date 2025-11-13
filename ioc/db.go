package ioc

import (
	"context"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongodb() *mongo.Client {
	type Config struct {
		Account  string `yaml:"account"`
		Address  string `yaml:"address"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
	}
	var config Config
	if err := viper.UnmarshalKey("mongo", &config); err != nil {
		panic(err)
	}
	// 拼接 URL 时补充 authSource
	url := "mongodb://" + config.Account + ":" + config.Password + "@" + config.Address + ":" + config.Port
	ClientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.Background(), ClientOptions)
	if err != nil {
		panic(err)
	}
	return client
}
