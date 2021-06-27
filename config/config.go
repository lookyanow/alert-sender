package config

import (
	"log"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config structure contains of app params
type Config struct {
	ServerHost  string
	ServerPort  string
	SmsURL      string
	SmsUser     string
	SmsPassword string
	Phones      string
	PhoneFile   string
}

// NewConfigFromEnv method get params from ENV
func NewConfigFromEnv() *Config {

	// default values
	pflag.String("server.host", "127.0.0.1", "host on which the server should listen")
	pflag.String("server.port", "8080", "port on which the server should listen")
	pflag.String("sms.url", "https://sms.city-srv.ru/sms", "sms gateurl")
	pflag.String("sms.user", "", "sms gateway user")
	pflag.String("sms.password", "", "sms gateway password")
	pflag.String("phones", "", "phone numbers list")
	pflag.String("phonefile", "", "phone list file in json format")

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatalf("Couldn't bind flags with error: %s", err)
	}

	// Environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	log.Printf("%v", viper.GetString("phones"))

	return &Config{
		ServerHost:  viper.GetString("server.host"),
		ServerPort:  viper.GetString("server.port"),
		SmsURL:      viper.GetString("sms.url"),
		SmsUser:     viper.GetString("sms.user"),
		SmsPassword: viper.GetString("sms.password"),
		Phones:      viper.GetString("phones"),
		PhoneFile:   viper.GetString("phonefile"),
	}
}
