package conf

import (
	"github.com/BurntSushi/toml"
	"log"
)

type database struct {
	Url string
}
type Config struct {
	DB        database `toml:"database"`
	JwtSecret string
}

var AppConfig Config

func SetUp() {
	if _, err := toml.DecodeFile("conf/app.toml", &AppConfig); err != nil {
		log.Fatal(err)
	}
}
