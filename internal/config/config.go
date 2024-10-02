package config

import (
	"github.com/kkyr/fig"
)

type Config struct {
	Token       string `fig:"Token"`
	WelcomeMsg  string `fig:"WelcomeMsg" default:"Привет, пришли мне бинарь\\картинку\\любой файл и я сделаю на его основе изображение."`
	CompliteDir string `fig:"CompliteDir"`
	BinsDir     string `fig:"BinsDir"`
}

func LoadConfig(path string) (*Config, error) {
	conf := new(Config)
	if err := fig.Load(conf, fig.File(path)); err != nil {
		return nil, err
	}
	

	return conf, nil
}
