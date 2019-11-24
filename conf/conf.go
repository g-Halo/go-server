package conf

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Env         string `json:"env"`
	TcpAddress  string `json:"tcp_address"`
	HttpAddress string `json:"http_address"`

	LogicRPCAddress string `json:"logic_rpc_address"`
	AuthRPCAddress  string `json:"auth_rpc_address"`
	HttpApiAddress  string `json:"http_api_address"`
	MongoDbAddress  string `json:"mongodb_address"`
	SecretKey       string `json:"secret_key"`
	DB              string `json:"db"`

	ChannelBucketCount int `json:"channel_bucket_count"`
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

func No_db() bool {
	return Conf.DB == ""
}

// 无 db，可以不指定 db，那么所有数据就存储在内存
func (c *Config) No_db() bool {
	return c.DB == ""
}
