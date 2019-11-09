# Auth

鉴权的模块

1. 颁发令牌
Params: username, password

流程

请求 logic 层校验 username,password 是否合法
若合法，则颁发令牌

Response: token

2. 校验令牌
Params: token
Response: ok | fail

