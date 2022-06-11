<!--
 * @Descripttion:
 * @Author: weihaoyu
-->

# gin-api
基于 Gin 开发的微服务框架，封装各种常用组件，目的在于提高 Go 应用开发效率。
<br><br>
如果您对本框架有任何意见或建议，欢迎随时通过以下方式反馈和完善：
1. 提 issues 反馈
2. 通过下方的联系方式直接联系我
3. 提 PR 共同维护
<br><br>

## 联系我
QQ群：909211071
<br>
个人QQ：444216978
<br>
微信：AirGo___
<br><br>

## 目前已支持
✅ &nbsp;多格式配置读取
<br>
✅ &nbsp;服务优雅关闭
<br>
✅ &nbsp;进程结束资源自动回收
<br>
✅ &nbsp;日志抽象和标准字段统一（请求、DB、Redis、RPC）
<br>
✅ &nbsp;DB
<br>
✅ &nbsp;RabbitMQ
<br>
✅ &nbsp;Redis
<br>
✅ &nbsp;分布式缓存（解决缓存穿透、击穿、雪崩）
<br>
✅ &nbsp;分布式链路追踪
<br>
✅ &nbsp;分布式锁
<br>
✅ &nbsp;服务注册
<br>
✅ &nbsp;服务发现
<br>
✅ &nbsp;负载均衡
<br>
✅ &nbsp;HTTP-RPC 超时传递
<br>
✅ &nbsp;端口多路复用
<br>
✅ &nbsp;gRPC
<br><br>

## 后续逐渐支持
日志收集
<br>
监控告警
<br>
限流
<br>
熔断
<br><br>

# 目录结构
```
- gin-api 
  - app //用户应用目录
    - conf //服务配置文件目录
      - dev
      - liantiao
      - online
      - qa
    - config //启动加载配置目录
      - app.go //应用配置
    - loader //资源加载
    - resource
      - resource.go //全局资源
    - response
      - response.go //http响应
    - router
      router.go //路由定义和中间件注册
    - rpc //三方rpc调用封装
      - gin-api //gin-api服务
    - module //各模块核心实现，按照业务边界划分目录
      - module1 //模块1
        - api //对外暴露api
        - job //离线任务
        - responsitory //存储层
        - service //核心业务代码
      - module1 //模块2
        - api //对外暴露api
        - job //离线任务
        - responsitory //存储层
        - service //核心业务代码
    - main.go //app入口文件
  - bootstrap //应用启动
  - client
    - codec //编码
    - grpc //grpc客户端
    - http //http客户端
  - server
    - grpc //grpc服务端
    - http //http服务端
  - library //基础组件库，不建议修改
    - apollo //阿波罗
    - cache //分布式缓存
    - config //配置加载
    - cron //任务调度
    - endless //endless
    - etcd //etcd
    - grpc //grpc封装
    - jaeger //jaeger分布式链路追踪
    - job //离线任务
    - lock //分布式锁
    - logger //日志
    - orm //db orm
    - rabbitmq //rabbitmq
    - redis //redis
    - registry //注册中心
    - selector //负载均衡器
    - servicer //下游服务
  .gitignore
  Dockerfile
  LICENSE
  Makefile
  README.md
  go.mod
  go.sum
```

# 配置相关
基于三方组件viper，文件配置需放到main.go同级目录conf/xxx下
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
panic 接口：http://localhost:777/test/panic
<br>
db 和 redis测试接口：http://localhost:777/test/conn （依赖 mysql 和 redis）
<br>
分布式链路追踪+服务注册+服务发现接口：http://localhost:777/test/rpc (依赖 mysql 和 redis）
<br><br>


# 分布式链路追踪
<img src="https://github.com/why444216978/images/blob/master/jaeger.png" width="800" height="300" alt="jaeger"/>
<br>

# 运行
1. app 目录下的资源初始化和服务注册是留给开发者自己扩展的，您可自行调整资源加载。
2. 查看 conf/xxx 目录的各个配置文件，改成符合自己的
3. log 配置中的目录确保本地存在且有写入权限
4. go run main.go -env dev（不带 -env 参数默认用 dev 配置）
<br>


**注意：测试 /test/conn 和 /test/rpc 接口时，应完成如下几项：**
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

## 运行HTTP服务
```
[why@bogon] ~/Desktop/go/gin-api/app$go run main.go -env dev -server http
2022/06/12 04:45:13 Actual pid is 5227
2022/06/12 04:45:13 The environment is :dev
2022/06/12 04:45:14 start http, port 8777
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /ping                     --> github.com/why444216978/gin-api/app/module/ping/api.Ping (4 handlers)
[GIN-debug] GET    /ping/rpc                 --> github.com/why444216978/gin-api/app/module/ping/api.RPC (4 handlers)
[GIN-debug] POST   /test/rpc                 --> github.com/why444216978/gin-api/app/module/test/api.Rpc (4 handlers)
[GIN-debug] POST   /test/rpc1                --> github.com/why444216978/gin-api/app/module/test/api.Rpc1 (4 handlers)
[GIN-debug] POST   /test/panic               --> github.com/why444216978/gin-api/app/module/test/api.Panic (4 handlers)
[GIN-debug] POST   /test/conn                --> github.com/why444216978/gin-api/app/module/goods/api.Do (4 handlers)
watching prefix:gin-api now...
service gin-api  put key: gin-api.192.168.1.104.777 val: {"Host":"192.168.1.104","Port":777}


[why@bogon] ~/Desktop$curl http://localhost:8777/ping
{"code":0,"toast":"success","data":{},"errmsg":"success","trace_id":""}
```

## 运行grpc服务
```
[why@bogon] ~/Desktop/go/gin-api/app$go run main.go -env dev -server grpc
2022/06/12 04:41:32 Actual pid is 3001
2022/06/12 04:41:32 The environment is :dev
2022/06/12 04:41:32 start grpc, port 8777

[why@bogon] ~/Desktop$curl http://localhost:8888/v1/example/echo
{"message":" world"}

[why@bogon] ~/Desktop/go/gin-api/app$go run main.go -job grpc-test
2022/06/12 04:59:04 Actual pid is 12538
2022/06/12 04:59:04 start job by grpc-test
message:"why world"
message:"why world"
message:"why world"
```



