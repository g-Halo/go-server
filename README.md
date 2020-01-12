### go-server
使用 go 实现的聊天室

### 部署
```
在本机搭建 golang 环境

git clone 本仓库

依赖安装
`go get ./...` 
需要注意代理问题.. 代理问题可参考 https://github.com/goproxy/goproxy.cn

配置文件
`cp config.json.example config.json`

运行项目
`go run main.go`

```

### API
#### 登录
POST /v1/login

#### 注册
POST /v1/sign

#### 获取联系人列表
GET /v1/contacts

#### WebSocket连接
GET /v1/ws

### TODO
1. 保证消息的一致性
2. 消息重复机制
3. ws连接状态的维护
4. 在线人数和聊天室的数量统计
5. 接入 grpc 
6. 各模块拆分解耦
7. 心跳检测的完善
9. 接入 db/redis 存储数据
8. 构建 docker 环境