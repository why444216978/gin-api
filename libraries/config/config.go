package config

import (
	"fmt"
	"gin-api/app_const"
	"gin-api/libraries/apollo"
	"log"

	"github.com/why444216978/go-util/conversion"
	util_file "github.com/why444216978/go-util/file"
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
	cfgList map[string]interface{}
)

func init() {
	cfgList = make(map[string]interface{}, app_const.CONFIGS_NUM)
}

func GetConfigToJson(file, section string) map[string]interface{} {
	ret := make(map[string]interface{}, 10)

	if app_const.CONFIG_SOURCE == SOURCE_APOLLO {
		cfg, _ := apollo.LoadApolloConf(app_const.SERVICE_NAME, []string{"application"})
		cfgMap, _ := conversion.JsonToMap(cfg[file])
		ret = cfgMap[section].(map[string]interface{})
	} else if app_const.CONFIG_SOURCE == SOURCE_JSON {
		ret = getJsonConfig(file, section)
	} else if app_const.CONFIG_SOURCE == SOURCE_INI {
		ret, _ = getIniConfig(file, section)
	}

	return ret
}

func getJsonConfig(file, section string) map[string]interface{} {
	if cfgList[file] != nil {
		return cfgList[file].(map[string]interface{})
	}

	jsonStr, _ := util_file.ReadJsonFile(path + file + ".json")
	cfgMap, _ := conversion.JsonToMap(jsonStr)
	cfgList[file] = cfgMap[section].(map[string]interface{})

	log.Println(fmt.Sprintf("load %s.json", file))

	return cfgList[file].(map[string]interface{})
}

func getIniConfig(file string, cfgSection string) (map[string]interface{}, error) {
	if cfgList[file] != nil {
		return cfgList[file].(map[string]interface{}), nil
	}

	configFile := fmt.Sprintf("%s%s.ini", path, file)
	iniFile, err := ini.Load(configFile)
	if err != nil {
		return nil, err
	}
	section := iniFile.Section(cfgSection)

	ret := make(map[string]interface{})
	cfgFields := section.KeyStrings()
	length := len(cfgFields)
	for i := 0; i < length; i++ {
		ret[cfgFields[i]] = section.Key(cfgFields[i]).String()
	}
	cfgList[file] = ret

	log.Println(fmt.Sprintf("load %s.json", file))

	return cfgList[file].(map[string]interface{}), err
}
