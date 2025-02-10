### tss (TIM-SupportServices) 
##### tss是tim 后端业务接口，用于服务器端与tim的信息交互，tssgo 是 tss的go实现

### 快速使用

```bash
go get github.com/donnie4w/tssgo
```

#### 连接tim的tcp adm服务
```go
client, err := tssgo.NewTssClient(false, "127.0.0.1:40000", "admin", "123", "")
if err != nil {
    panic(err)
}
```

##### 参数说明： tls bool, addr string, name, pwd,domain string
* tls  使用tls安全协议连接
* addr 服务器地址与端口
* name 后台登录用户名
* pwd  登录密码
* domain 域名(选填)




