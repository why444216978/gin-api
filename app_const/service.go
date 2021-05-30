package app_const

const (
	SERVICE_NAME  = "gin-api"
	SERVICE_PORT  = 777
	PRODUCT       = "gin-api"
	MODULE        = "gin-api"
	ENV           = "development"
	CONFIG_SOURCE = "ini" //apollo、json、ini
	CONFIGS_NUM   = 10    //配置文件数，影响配置file享元map初始化大小
	DB_NUM        = 10    //数据库实例数，影响数据库连接享元map初始化大小
)
