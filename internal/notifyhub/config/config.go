package config

import (
	"flag"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type ApplicationConfig struct {
	Server   ServerConfig
	Database DatabaseConfig
	Broker   BrokerConfig
	Email    EmailConfig
	Telegram TelegramConfig
	Sms      SmsConfig
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

type BrokerConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

type EmailConfig struct {
	Host               string
	Port               int
	Login              string
	Password           string
	UseTLS             bool
	InsecureSkipVerify bool
	TLSCrtFile         string
	TLSKeyFile         string
}

type TelegramConfig struct {
	BotToken string
}

type SmsConfig struct {
	Login string
	Key   string
}

func setDefaults(cfg *viper.Viper) {

	cfg.SetDefault("Database.Host", "localhost")
	cfg.SetDefault("Database.Port", 5432)
	cfg.SetDefault("Database.User", "notifyhub_user")
	cfg.SetDefault("Database.Password", "1")
	cfg.SetDefault("Database.Database", "notifyhub")
	cfg.SetDefault("Database.SSLMode", "disable")

	cfg.SetDefault("Broker.Host", "localhost")
	cfg.SetDefault("Broker.Port", 5672)
	cfg.SetDefault("Broker.User", "guest")
	cfg.SetDefault("Broker.Password", "guest")
}

func ReadConfig(configFilePath string) (*ApplicationConfig, error) {

	appConfig := &ApplicationConfig{}

	cfg := viper.New()

	cfg.SetEnvPrefix("notifyhub")
	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cfg.AutomaticEnv()

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	_ = cfg.BindPFlags(pflag.CommandLine)

	setDefaults(cfg)

	cfg.SetConfigFile(configFilePath)

	err := cfg.ReadInConfig()
	if err != nil {
		return nil, err
	}

	if err := cfg.Unmarshal(appConfig); err != nil {
		return nil, err
	}

	return appConfig, nil
}
