<!--
 * @Descripttion:
 * @Author: weihaoyu
-->

# gin-api
基于 Gin 开发的微服务框架，封装各种常用组件，目的在于提高Go应用开发效率。
<br><br>

## 目前已支持
✅ &nbsp;多格式配置读取
<br>
✅ &nbsp;日志抽象和标准统一（请求、DB、Redis、RPC）
<br>
✅ &nbsp;DB
<br>
✅ &nbsp;RabbitMQ
<br>
✅ &nbsp;Redis
<br>
✅ &nbsp;缓存处理（解决缓存穿透、击穿、雪崩）
<br>
✅ &nbsp;服务注册
<br>
✅ &nbsp;服务发现
<br>
✅ &nbsp;分布式锁
<br>
✅ &nbsp;分布式链路追踪
<br><br>

## 后续逐渐支持
<br>
gRPC
<br>
端口多路复用
<br>
日志收集
<br>
监控告警
<br>
限流
<br>
熔断
<br><br>

# 配置相关
基于三方组件viper，文件配置需放到main.go同级目录conf_xx下
<br><br>

# 统一日志
基于 zap 二次封装，抽象统一接口、数据库、缓存、RPC 调用日志结构，便于后期日志收集和搜索
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
<br><br>


# 分布式链路追踪
<img src="https://github.com/why444216978/images/blob/master/jaeger.png" width="800" height="300" alt="jaeger"/>
<br>

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
[why@localhost] ~/Desktop/go/gin-api$go run main.go -env dev
Actual pid is 34592
The environment is :dev
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /ping                     --> github.com/why444216978/gin-api/api/ping.Ping (5 handlers)
[GIN-debug] POST   /test/rpc                 --> github.com/why444216978/gin-api/api/test.Rpc (5 handlers)
[GIN-debug] POST   /test/rpc1                --> github.com/why444216978/gin-api/api/test.Rpc1 (5 handlers)
[GIN-debug] POST   /test/panic               --> github.com/why444216978/gin-api/api/test.Panic (5 handlers)
[GIN-debug] POST   /test/conn                --> github.com/why444216978/gin-api/api/conn.Do (5 handlers)
start by server
Start with :777
watching prefix:gin-api now...
service gin-api  put key: gin-api.192.168.1.104.777 val: {"Host":"192.168.1.104","Port":777}
```

# 联系我
QQ群：909211071
<br>
个人QQ：444216978
<br>
微信：AbleYu_
<br>
有疑问随时联系本人
<br>
