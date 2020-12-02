package config

import (
	"fmt"
	"gin-api/app_const"
	"gin-api/libraries/apollo"
	"gin-api/libraries/util/conversion"
	util_file "gin-api/libraries/util/file"
	"gopkg.in/ini.v1"
)

type Config struct {
	Result map[string]string
	Err    string
}

const (
	SOURCE_APOLLO = "apollo"
	SOURCE_JSON   = "json"
	SOURCE_INI    = "ini"
)

const path = "./configs/"

var (
	cfgList map[string]*ini.File
)

func init() {
	cfgList = make(map[string]*ini.File, app_const.CONFIGS_NUM)
}

func GetConfigToJson(file, section string) map[string]interface{} {
	ret := make(map[string]interface{}, 10)

	if app_const.CONFIG_SOURCE == SOURCE_APOLLO {
		cfg := apollo.LoadApolloConf(app_const.SERVICE_NAME, []string{"application"})
		cfgMap := conversion.JsonToMap(cfg[file])
		ret = cfgMap[section].(map[string]interface{})
	} else if app_const.CONFIG_SOURCE == SOURCE_JSON {
		return getJsonConfig(file, section)
	} else if app_const.CONFIG_SOURCE == SOURCE_INI {
		return getIniConfig(file, section)
	} else {
		panic("log source type error")
	}
	return ret
}

func getJsonConfig(file, section string) map[string]interface{}{
	jsonStr := util_file.ReadJsonFile(path + file + ".json")
	cfgMap := conversion.JsonToMap(jsonStr)
	return cfgMap[section].(map[string]interface{})
}

func getIniConfig(cfgType string, cfgSection string) map[string]interface{} {
	ret := make(map[string]interface{})

	configFile := fmt.Sprintf("%s%s.ini", path, cfgType)
	file,err := ini.Load(configFile)
	if err != nil{
		panic(err)
	}
	section := file.Section(cfgSection)

	cfgFields := section.KeyStrings()
	length := len(cfgFields)
	for i := 0; i < length; i++ {
		ret[cfgFields[i]] = section.Key(cfgFields[i]).String()
	}

	return ret
}