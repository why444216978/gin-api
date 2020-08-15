package config

func GetLogFields() map[string]string {
	logFields := make(map[string]string, 3)
	logFieldsSection := "log_fields"
	logFieldsConfig := GetConfig("log", logFieldsSection)
	logFields["query_id"] = logFieldsConfig.Key("query_id").String()
	logFields["header_id"] = logFieldsConfig.Key("header_id").String()
	logFields["header_hop"] = logFieldsConfig.Key("header_hop").String()
	return logFields
}

func GetXhopField() string {
	logFieldsSection := "log_fields"
	logFieldsConfig := GetConfig("log", logFieldsSection)
	return logFieldsConfig.Key("header_hop").String()
}
