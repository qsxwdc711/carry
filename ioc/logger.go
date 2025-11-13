package ioc

import (
	"context"
	"encoding/json"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"runtime"
	"sx-go/internal/domain"
	"time"
)

type mongoWriter struct {
	logColl *mongo.Collection
}

func NewMongoWriter(db *mongo.Client) mongoWriter {
	databaseName := viper.GetString("mongo.database")
	return mongoWriter{
		logColl: db.Database(databaseName).Collection("log"),
	}
}

func InitLogger(db *mongo.Client) domain.Loggers {
	// 创建 encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	console := zapcore.Lock(os.Stdout)
	var core zapcore.Core
	mongodb := NewMongoWriter(db)
	if runtime.GOOS == "linux" {
		//linux下才会计入到数据库中
		core = zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(&mongodb), zap.NewAtomicLevel())
	} else {
		//core = zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(&mongodb), zap.NewAtomicLevel())
		core = zapcore.NewCore(encoder, console, zap.NewAtomicLevel())
	}
	logInfo := zap.New(zapcore.NewCore(encoder, console, zap.NewAtomicLevel()), zap.AddCaller())
	// 创建 logger
	log := zap.New(core, zap.AddCaller())
	defer log.Sync()
	zap.ReplaceGlobals(log)
	//使用全局log，也就是logg会计入到mongo中，
	//loginfo输出到控制台
	return domain.Loggers{
		Logg:    log,
		LogInfo: logInfo,
	}
}

func (mw mongoWriter) Write(p []byte) (n int, err error) {
	var logMap map[string]interface{}
	if err = json.Unmarshal(p, &logMap); err != nil {
		return 0, err
	}
	var log logStruct
	if logMap["caller"] != nil {
		log.Caller = logMap["caller"].(string)
	}
	if logMap["error"] != nil {
		log.Error = logMap["error"].(string)
	}
	if logMap["route"] != nil {
		log.Route = logMap["route"].(string)
	}
	if logMap["msg"] != nil {
		log.Msg = logMap["msg"].(string)
	}
	log.Level = logMap["level"].(string)
	log.Time = time.Now().Format("2006-01-02 15:04:05")
	_, err = mw.logColl.InsertOne(context.Background(), log)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

type logStruct struct {
	Level  string `bson:"level" json:"level"`
	Time   string `bson:"time" json:"ts"`
	Caller string `bson:"caller" json:"caller"`
	Msg    string `bson:"msg" json:"msg"`
	Route  string `bson:"route" json:"route"`
	Error  string `bson:"error" json:"error"`
}
