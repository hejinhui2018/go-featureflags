# POST 重定向复现

## 可观察行为

客户端提交 POST 后跟随服务端重定向，最终到达规范地址时请求方法变成 GET，原始正文也丢失。GET 请求仍然能够到达规范地址。

## 复现命令

```powershell
go test ./... -run '^TestPostRedirectPreservesMethodAndBody$' -count=20
```

基线中目标测试每次都会看到最终方法为 GET 且正文为空。期望客户端跟随重定向后仍以 POST 发送原始正文。
