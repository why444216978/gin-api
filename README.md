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
支持apollo、json、ini三种格式，文件配置需放到main.go同级目录configs下
<br>
- apollo：嵌套json格式，用于兼容mysql、redis等多实例
- json：嵌套json格式，用于兼容mysql、redis等多实例
- ini：section格式，用于兼容mysql、redis等多实例
<br>
通过 app_const.server.CONFIG_SOURCE 变量切换
<br>

```
package app_const

const (
	SERVICE_NAME  = "purchase-server"
	SERVICE_PORT  = 777
	PRODUCT       = "gin-api"
	MODULE        = "gin-api"
	ENV           = "development"
	CONFIG_SOURCE = "ini" //apollo、json、ini
	CONFIGS_NUM   = 10     //配置文件数，影响配置file享元map初始化大小
)

```

## ini（默认格式）

### env.ini example:

```
[env]
env = development
```

### log.ini example：
```
[log]
dir = /data/logs
area = 1
query_field = "logid"
header_field = "X-Logid"
```


### mysql.ini example:

```
[default_read]
host = 127.0.0.1
user = why
password = why123
port = 3306
db = why
charset = utf8
max_open = 8
max_idle = 4
exec_timeout = 10000

[default_write]
host = 127.0.0.1
user = why
password = why123
port = 3306
db = why
charset = utf8
max_open = 8
max_idle = 4
exec_timeout = 10000
```

### redis.ini example:

```
[default]
host = 127.0.0.1
port = 6379
db = 0
auth =
max_active = 600
max_idle = 10
exec_timeout = 100000
```

### log.ini example:

```
[log]
dir = /data/logs/
area = 1
query_field = "logid"
header_field = "X-Logid"
```

### es.ini example:

```
[default]
host = http://127.0.0.1
port = 9200
```

### env.json example：
```
{
  "env":{
    "env": "development"
  }
}
```

### log.json example：
```
{
  "log":{
    "dir":"/data/logs",
    "area":1,
    "query_field":"logid",
    "header_field":"X-Logid"
  }
}
```

# 运行

1. 创建上述基础配置文件
2. log.ini中的dir目录确保本地存在且有写入权限
3. go run main.go
4. curl localhost:777/ping

