package conf

import (
	"github.com/BurntSushi/toml"
	"log"
)

type database struct {
	Url string
}
type path struct {
	Avatar   string
	File     string
	Video    string
	Frontend string
	Cover    string
}
type Config struct {
	JwtSecret string
	DB        database `toml:"database"`
	Path      path
}

var AppConfig Config

func SetUp() {
	if _, err := toml.DecodeFile("conf/app.toml", &AppConfig); err != nil {
		log.Fatal(err)
	}
}
