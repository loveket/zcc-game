package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"xiuianserver/utils"
)

var GlobalConfig *YamlConfig

func init() {
	GlobalConfig = &YamlConfig{}
	data, err := os.ReadFile(utils.GetOsPwd() + "\\config\\global.yml")
	if err != nil {
		fmt.Println("read global err", err)
		return
	}
	err = yaml.Unmarshal(data, GlobalConfig)
	if err != nil {
		fmt.Println("read global yaml err", err)
		return
	}
}

type YamlConfig struct {
	HttpConfig    HTTPConfig    `yaml:"Http"`
	NsqConfig     NSQConfig     `yaml:"Nsq"`
	RedisConfig   REDISConfig   `yaml:"Redis"`
	ElasticConfig ELASTICConfig `yaml:"Elastic"`
	RpcConfig     RPCConfig     `yaml:"Rpc"`
}

type HTTPConfig struct {
	IpAddr string `yaml:"IpAddr"`
}
type NSQConfig struct {
	RemoteAddr string `yaml:"RemoteAddr"`
}
type RPCConfig struct {
	IpAddr string `yaml:"IpAddr"`
}
type REDISConfig struct {
	RemoteAddr string `yaml:"RemoteAddr"`
	User       string `yaml:"User"`
	Pass       string `yaml:"Pass"`
	DB         int    `yaml:"DB"`
}
type ELASTICConfig struct {
	HttpAddr string `yaml:"HttpAddr"`
	NodeAddr string `yaml:"NodeAddr"`
	Sniff    bool   `yaml:"Sniff"`
}
