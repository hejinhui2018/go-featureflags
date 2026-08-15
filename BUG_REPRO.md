# 反向代理逐跳请求头复现

## 可观察行为

客户端通过 `Connection` 声明仅用于当前连接的自定义请求头后，准备转发的后端请求仍包含这些自定义头。标准逐跳头会被删除，普通端到端请求头可以正常保留。

## 复现命令

```powershell
go test ./... -run '^TestPrepareBackendRequestRemovesConnectionNominatedHeaders$' -count=20
```

基线中每次运行都会发现 `X-Internal-Hop` 和 `X-Debug-Hop` 被转发。期望后端请求删除这些逐跳头，同时不修改原始请求，也不删除普通端到端请求头。
