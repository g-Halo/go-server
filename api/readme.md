### Api
> gRpc 的所有接口方法的实现

#### Auth: 实现身份校验的接口
- SignIn 登录
- SignUp 注册
- Validate 校验 token 有效性

#### Logic: 逻辑层接口
- PushMessage 发送消息
- GetUser 获取某用户的详情
- GetUsers 获取用户列表
- FindOrCreateRoom 房间的查询（如查询不到就新建）