package conf

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)

type Config struct {
	Env         string `json:"env"`
	TcpAddress  string `json:"tcp_address"`
	HttpAddress string `json:"http_address"`

	CommetAddress    string `json:"commet_address"`
	LogicRPCAddress  string `json:"logic_rpc_address"`
	AuthRPCAddress   string `json:"auth_rpc_address"`
	HttpApiAddress   string `json:"http_api_address"`
	WebSocketAddress string `json:"websocket_address"`
	MongoDbAddress   string `json:"mongodb_address"`
	SecretKey        string `json:"secret_key"`
	DB               string `json:"db"`

	NSQAddress string `json:"nsq_address"`
	NSQTopic string `json:"nsq_topic"`

	RoomChannelsCount int `json:"room_channels_count"`
}

var Conf *Config
var confPath string

func init() {
	flag.StringVar(&confPath, "conf", "config.json", "default config path")
}

func LoadConf() *Config {
	if Conf == nil {
		bytes, err := ioutil.ReadFile(confPath)
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
