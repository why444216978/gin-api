<!--
 * @Descripttion:
 * @Author: weihaoyu
-->

# gin-api

基于 go-gin 开发的 api 框架，封装各种常用组件
<br>
有疑问随时联系我
<br>
QQ：444216978
<br>
微信：AbleYu_
<br>

# run

```
go run main.go

curl localhost:777/test/ping
```

# 配置相关
配置放到main.go同级目录configs下

# app.ini example:

```
[app]
env = development
port = 777
app_id = gin-api
product = gin-api
module = gin-api
```

# mysql.ini example:

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

# redis.ini example:

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

# log.ini example:

```
[run]
dir = ./logs/run/
area = 1

[error]
dir = ./logs/error/
area = 1

[amqp]
dir = ./logs/amqp/
area = 1
```

# es.ini example:

```
[default]
host = http://127.0.0.1
port = 9200
```
