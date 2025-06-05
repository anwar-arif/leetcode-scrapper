package config

import (
	"github.com/spf13/viper"
	"log"
	"sync"
)

type Headers struct {
	XCsrfToken string `yaml:"x-csrftoken" json:"-"`
}

func readConfig(filename string) {
	viper.SetConfigFile(filename)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("unable to read config file")
	}
}

type Application struct {
	Headers Headers `yaml:"headers" json:"-"`
}

func loadApp(filename string) error {
	readConfig(filename)
	viper.AutomaticEnv()

	appConfig = &Application{
		Headers: Headers{XCsrfToken: viper.GetString("app.secrets.csrf_token")},
	}
	log.Println("app config loaded")

	return nil
}

var appOnce = sync.Once{}
var appConfig *Application

func GetApp(filename string) *Application {
	appOnce.Do(func() {
		loadApp(filename)
	})
	return appConfig
}
