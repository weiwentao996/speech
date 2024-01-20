package config

import "github.com/spf13/viper"

type Config struct {
	App struct {
		Port int `mapstructure:"port"`
	} `json:"app"`
	Azure struct {
		Key          string `mapstructure:"key"`
		Region       string `mapstructure:"region"`
		DefaultVoice string `mapstructure:"default_voice"`
	}
}

var AppConfig Config

func init() {
	viper.AddConfigPath("./")     //设置读取的文件路径
	viper.SetConfigName("config") //设置读取的文件名
	viper.SetConfigType("yaml")   //设置文件的类型
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&AppConfig); err != nil {
		panic(err)
	}
}
