<!--
 * @Descripttion:
 * @Author: weihaoyu
-->
<img src="https://github.com/why444216978/images/blob/master/qrcode.png" width="300" height="100" alt="公众号"/>

# gin-api
基于 Gin 开发的 api 框架，封装各种常用组件，包括配置、日志、DB、RabbitMQ、Redis、缓存处理（解决缓存穿透、击穿、雪崩）、分布式链路追踪等，目的在于提高Go应用开发效率。
<br><br>


# 配置相关
基于三方组件viper，文件配置需放到main.go同级目录conf_xx下
<br><br>

# 服务发现
目前支持 etcd， <a href="https://success.blog.csdn.net/article/details/119827014">集群搭建教程</a>，相关配置无需更改的情况下，按照教程搭建运行即可。
<br><br>

# 测试
检测接口：http://localhost:777/ping 
<br>
panic接口：http://localhost:777/test/panic
<br>
db和redis测试接口：http://localhost:777/test/conn
<br>
分布式链路追踪+服务注册+服务发现接口：http://localhost:777/test/rpc
<br>
服务发现测试脚本：go run main.go -job registry
<br><br>


# 分布式链路追踪
<img src="https://github.com/why444216978/images/blob/master/jaeger.png" width="800" height="300" alt="jaeger"/>
<br><br>

# 运行
1. 查看 conf_dev 目录的各个配置文件，改成符合自己的
2. log 配置中的目录确保本地存在且有写入权限
3. go run main.go
<br>


**注意：测试 /test/conn 和 /test/rpc 接口时，应确检查如下几项：**
1. 创建 test 库
2. 创建 test 表并随意插入数据
```
CREATE TABLE `test` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `goods_id` bigint(20) unsigned NOT NULL,
  `name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `idx_goods` (`goods_id`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin 
```

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
