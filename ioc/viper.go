package ioc

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func InitViper() {
	cfile := pflag.String("config", "config/conf.yaml", "指定配置文件路径")
	//顺序不能变可以在启动的时候传参数
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	//开启监视
	viper.WatchConfig()
	//回调
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//	//配置文件修改会进入这里
	//	//fmt.Println(in.Name, in.Op)
	//})
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
