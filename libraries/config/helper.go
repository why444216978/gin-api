package config

func GetLogFields() map[string]string {
	logFields := make(map[string]string, 3)
	logFieldsSection := "log_fields"
	logFieldsConfig := GetConfig("log", logFieldsSection)

	logFields["query_id"] = logFieldsConfig.Key("query_id").String()
	if logFields["query_id"] == "" {
		logFields["query_id"] = "logid"
	}

	logFields["header_id"] = logFieldsConfig.Key("header_id").String()
	if logFields["header_id"] == "" {
		logFields["header_id"] = "LOG-ID"
	}

	logFields["header_hop"] = logFieldsConfig.Key("header_hop").String()
	if logFields["header_hop"] == "" {
		logFields["header_hop"] = "X-HOP"
	}

	return logFields
}

func GetHeaderLogIdField() string {
	logFieldsSection := "log_fields"
	logFieldsConfig := GetConfig("log", logFieldsSection)

	field := logFieldsConfig.Key("header_id").String()
	if field == "" {
		field = "LOG-ID"
	}

	return field
}

func GetLogIdField() string {
	logFieldsSection := "log_fields"
	logFieldsConfig := GetConfig("log", logFieldsSection)

	field := logFieldsConfig.Key("query_id").String()
	if field == "" {
		field = "logid"
	}

	return field
}

func GetXhopField() string {
	logFieldsSection := "log_fields"
	logFieldsConfig := GetConfig("log", logFieldsSection)

	field := logFieldsConfig.Key("header_hop").String()
	if field == "" {
		field = "X-HOP"
	}
	return field
}
