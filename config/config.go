package config

import (
	"github.com/golang/glog"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var (
	cnf *ConfigInfo
)

type ConfigInfo struct {
	Server *ServerConf `yaml:"Server" json:"Server"`
}

type ServerConf struct {
	Mode         string `yaml:"Mode"  json:"Mode"`
	Port         string `yaml:"Port" json:"Port"`
	ReadTimeout  int    `yaml:"ReadTimeout" json:"ReadTimeout"`
	WriteTimeout int    `yaml:"WriteTimeout" json:"WriteTimeout"`
}

func Config() *ConfigInfo {
	if cnf == nil {
		panic("cnf == nil")
	}
	return cnf
}

func createConfigByLocal(filename string) {
	glog.Infoln("读取本地配置文件")
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		glog.Error(err)
		return
	}
	err = yaml.Unmarshal(data, &cnf)
	if err != nil {
		glog.Error(err)
		return
	}

}

func InitConfig(configUrl string) {
	createConfigByLocal(configUrl)

}
