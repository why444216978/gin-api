package config

import (
	"gin-api/configs"
	"gin-api/libraries/apollo"
	"gin-api/libraries/util/conversion"
)

const (
	SOURCE_APOLLO = "apollo"
	SOURCE_FILE   = "file"
)

var (
	queryLogIdField  string
	headerLogIdField string
	logDir           string
	logArea          int
)

func GetLogConfig(source string) (string, int) {
	if logDir != "" && logArea != 0 {
		return logDir, logArea
	}

	if source == SOURCE_APOLLO {
		//apollo获取
		cfg := apollo.LoadApolloConf(configs.SERVICE_NAME, []string{"application"})
		logCfg := conversion.JsonToMap(cfg["log"])
		runLogDir := logCfg["dir"]
		tmpLogArea, _ := logCfg["area"]
		logArea := int(tmpLogArea.(float64))
		logDir := runLogDir.(string) + "/" + configs.SERVICE_NAME
		return logDir, logArea
	} else if source == SOURCE_FILE {
		errLogSection := "log"
		errorLogConfig := GetConfig("log", errLogSection)
		logDir := errorLogConfig.Key("dir").String()
		logDir = logDir + "/" + configs.SERVICE_NAME
		logArea, err := errorLogConfig.Key("area").Int()
		if err != nil {
			panic(err)
		}
		return logDir, logArea
	} else {
		panic("log source type error")
	}
}

func GetHeaderLogIdField(source string) string {
	if headerLogIdField != "" {
		return headerLogIdField
	}

	if source == SOURCE_APOLLO {
		cfg := apollo.LoadApolloConf(configs.SERVICE_NAME, []string{"application"})
		logCfg := conversion.JsonToMap(cfg["log"])
		headerLogIdField = logCfg["header_field"].(string)
	} else if source == SOURCE_FILE {
		logFieldsConfig := GetConfig("log", "log")

		headerLogIdField = logFieldsConfig.Key("header_field").String()
	} else {
		panic("log source type error")
	}

	if headerLogIdField == "" {
		headerLogIdField = "X-Logid"
	}

	return headerLogIdField
}

func GetQueryLogIdField(source string) string {
	if queryLogIdField != "" {
		return queryLogIdField
	}

	if source == SOURCE_APOLLO {
		cfg := apollo.LoadApolloConf(configs.SERVICE_NAME, []string{"application"})
		logCfg := conversion.JsonToMap(cfg["log"])
		queryLogIdField = logCfg["query_field"].(string)
	} else if source == SOURCE_FILE {
		logFieldsConfig := GetConfig("log", "log")

		queryLogIdField = logFieldsConfig.Key("query_field").String()
	} else {
		panic("log source type error")
	}

	if queryLogIdField == "" {
		queryLogIdField = "X-Logid"
	}

	return queryLogIdField
}
