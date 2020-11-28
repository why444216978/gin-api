package config

import (
	"fmt"
	"gin-api/app_const"
	"github.com/larspensjo/config"
	"gopkg.in/ini.v1"
	"log"
)

type Config struct {
	Result map[string]string
	Err    string
}

const path = "./configs/"

var (
	cfgList     map[string]*ini.File
)

func init() {
	cfgList = make(map[string]*ini.File, app_const.CONFIGS_NUM)
}

func GetConfig(cfgType string, cfgSection string) *ini.Section {
	if cfgList[cfgType] == nil {
		log.Println(fmt.Sprintf("load %s config file ", cfgType))
		var err error
		configFile := fmt.Sprintf("%s%s.ini", path, cfgType)
		cfgList[cfgType], err = ini.Load(configFile)
		if err != nil {
			panic(err)
		}
	}

	section := cfgList[cfgType].Section(cfgSection)

	return section
}

func (self *Config) getConfig(conn string, configFile string) {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	//flag.Parse()

	cfg, err := config.ReadDefault(configFile) //读取配置文件，并返回其Config

	if err != nil {
		logMsg := fmt.Sprintf("Fail to find %v,%v", configFile, err)
		self.Err = logMsg
	}

	self.Result = map[string]string{}
	if cfg.HasSection(conn) { //判断配置文件中是否有section（一级标签）
		options, err := cfg.SectionOptions(conn) //获取一级标签的所有子标签options（只有标签没有值）
		if err == nil {
			for _, v := range options {
				optionValue, err := cfg.String(conn, v) //根据一级标签section和option获取对应的值
				if err == nil {
					self.Result[v] = optionValue
				}
			}
		}
	}
}

func GetConfigEntrance(cfgType string, cfgSection string) map[string]string {
	cfg := new(Config)
	cfg.getConfig(cfgSection, path+cfgType+".ini")

	return cfg.Result
}
