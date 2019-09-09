package conf

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	TcpAddress 		string	`json:"tcp_address"`
	HttpAddress 	string	`json:"http_address"`
	MongoDbAddress 	string	`json:"mongodb_address"`
	SecretKey 		string 	`json:"secret_key"`
}

const configPath = "config.json"
var Conf *Config

func LoadConf() *Config {
	if Conf == nil {
		bytes, err := ioutil.ReadFile(configPath)
		if nil != err {
			panic(err)
		}

		conf := &Config{}
		if err = json.Unmarshal(bytes, &conf); err != nil {
			panic(err)
		}
		Conf = conf
	}

	return Conf
}