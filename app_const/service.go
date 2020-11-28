package app_const

const (
	SERVICE_NAME    = "gin-api"
	SERVICE_PORT    = 777
	PRODUCT         = "gin-api"
	MODULE          = "gin-api"
	ENV             = "development"
	CONFIG_SOURCE   = "file" //apollo or file
	CONFIGS_NUM     = 10     //配置文件数，影响配置file享元map初始化大小
	CONFIGS_SECTION = 10     //配置文件section数，影响配置section享元map初始化大小
)
