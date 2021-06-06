<!--
 * @Descripttion:
 * @Author: weihaoyu
-->

# gin-api

基于 go-gin 开发的 api 框架，封装各种常用组件
<br>
有疑问随时联系本人
<br>
QQ群：909211071
<br>
个人QQ：444216978
<br>
微信：AbleYu_
<br>

# 配置相关
基于三方组件viper，文件配置需放到main.go同级目录conf_xx下
<br>

```
package app_const

const (
	SERVICE_NAME  = "gin-api"
	SERVICE_PORT  = 777
	CONFIG_SOURCE = "ini" //apollo、json、ini
	CONFIGS_NUM   = 10    //配置文件数，影响配置file享元map初始化大小
)

```

## ini（默认格式）

### log.toml：
```
[log]
InfoFile = "./logs/info.log"
ErrorFile = "./logs/error.wf.log"
Level = "info"
```


### test_mysql:

```
[master]
Host = "127.0.0.1"
Port = "3306"
User = "root"
Password = "123456"
DB = "test"
Charset = "utf8mb4"
MaxOpen = 8
MaxIdle = 4
ExecTimeout = 10000

[slave]
Host = "127.0.0.1"
Port = "3306"
User = "root"
Password = "123456"
DB = "test"
Charset = "utf8mb4"
MaxOpen = 8
MaxIdle = 4
ExecTimeout = 10000
```

### default_redis.toml:

```
Host = "127.0.0.1"
Port = 6379
Auth = ""
DB = 0
ConnectTimeout = 1
ReadTimeout = 1
WriteTimeout = 1
MaxActive = 30
MaxIdle = 10
IsLog = true
ExecTimeout = 100000
```

### jaeger.toml:

```
Host = "127.0.0.1"
Port = "6831"
```


# 运行

1. 创建上述基础配置文件
2. log配置中的目录确保本地存在且有写入权限
3. go run main.go


```
[why@localhost] ~/Desktop/go/gin-api$go run main.go 
2020/12/20 17:44:43 load redis.json
2020/12/20 17:44:43 load mysql.json
2020/12/20 17:44:43 new test_dao
2020/12/20 17:44:43 new test_service
2020/12/20 17:44:43 load log.json
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

2020/12/20 17:44:43 load env.json
[GIN-debug] GET    /ping                     --> gin-api/controllers/ping.Ping (6 handlers)
[GIN-debug] GET    /test/rpc                 --> gin-api/controllers/opentracing.Rpc (6 handlers)
[GIN-debug] GET    /test/panic               --> gin-api/controllers/opentracing.Panic (6 handlers)
[GIN-debug] GET    /test/conn                --> gin-api/controllers/conn.Do (6 handlers)
2020/12/20 17:44:43 Actual pid is 9104
```


**注意：测试 /test/conn 接口时，应确检查如下几项：**
1. mysql 和 redis 配置文件符合示例配置文件中的default（当然可以自定义，不过需要更改 test_model.go 和 goods_service.go 中的 DB_NAME ）
2. 创建 test 库
3. 创建 test 表并随意插入数据
```
CREATE TABLE `test` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `goods_id` bigint(20) unsigned NOT NULL,
  `name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `idx_goods` (`goods_id`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin 
```


