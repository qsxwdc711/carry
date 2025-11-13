package domain

import "go.uber.org/zap"

type Loggers struct {
	Logg    *zap.Logger
	LogInfo *zap.Logger
}
type Log struct {
	Level  string `bson:"level" json:"level"`
	Time   string `bson:"time" json:"ts"`
	Caller string `bson:"caller" json:"caller"`
	Msg    string `bson:"msg" json:"msg"`
	Route  string `bson:"route" json:"route"`
	Error  string `bson:"error" json:"error"`
}
