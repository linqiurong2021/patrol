package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Conf struct {
	App      App      `yaml:"app"`
	Postgre  Postgre  `yaml:"postgre"`
	Redis    Redis    `yaml:"redis"`
	Register Register `yaml:"register"`
	Grpc     Grpc     `yaml:"grpc"`
}

// app配置信息
type App struct {
	Port string `yaml:"port"`
}

// 数据库配置信息
type Postgre struct {
	Host     string `yaml:"host"`
	Port     int32  `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

// redis 配置信息
type Redis struct {
	Host     string `yaml:"host"`
	Port     int32  `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Selected string `yaml:"selected"`
}

type Register struct {
	Host string `yaml:"host"`
}

type Grpc struct {
	Port string `yaml:"port"`
	Name string `yaml:"name"`
}

var Config Conf

// 初始化配置文件
func NewConf(path string) error {
	//
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("init conf error : %s",err)
	}
	return yaml.Unmarshal(bytes,&Config)
}