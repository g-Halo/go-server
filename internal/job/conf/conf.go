package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
)

type Config struct {
	MQ string `json:"mq"`
	Nsq *Nsq	`json:"nsq"`
}

type Nsq struct {
	Topic string `json:"topic"`
	Channel string `json:"channel"`
	LookUpAddress string `json:"nsqloopupd"`
}

var (
	confPath string
	Conf *Config
)

func init() {
	flag.StringVar(&confPath, "config", "job-example.toml", "default config path")
}

func Init() error {
	Conf = &Config{}
	if _, err := toml.DecodeFile(confPath, &Conf); err != nil {
		return err
	}

	return nil
}