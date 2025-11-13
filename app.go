package main

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	server  *gin.Engine
	mongodb *mongo.Client
	//redis   redis.Cmdable
	//cron    *cron.Cron
}
