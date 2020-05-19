package conf

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

type database struct {
	Url string
}
type email struct {
	Host     string
	Port     int
	Username string
	Password string
}
type path struct {
	Avatar   string
	File     string
	Video    string
	Frontend string
	Cover    string
}
type redis struct {
	Addr     string
	Password string
}
type Config struct {
	JwtSecret string
	Port      string
	Email     email
	DB        database `toml:"database"`
	Path      path
	Redis     redis
}

var AppConfig Config

func SetUp() {
	if _, err := toml.DecodeFile("conf/app.toml", &AppConfig); err != nil {
		log.Fatal(err)
	}
	InitDirectory()
}
func InitDirectory() { //如果目录不存在把目录建好
	_ = os.MkdirAll(AppConfig.Path.Video, os.ModePerm)
	_ = os.MkdirAll(AppConfig.Path.File, os.ModePerm)
	_ = os.MkdirAll(AppConfig.Path.Cover, os.ModePerm)
	_ = os.MkdirAll(AppConfig.Path.Avatar, os.ModePerm)
	//_=os.MkdirAll(AppConfig.Path.Frontend,os.ModePerm)
}
