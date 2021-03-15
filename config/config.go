package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Slack Slack
}

type Slack struct {
	ChannelId string `envconfig:"CHANNEL_ID"`
	Token     string `envconfig:"TOKEN"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
