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
