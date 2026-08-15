# 到期时刻仍返回旧值

## 问题描述

在 2026-08-15 12:00:00 存入 `session-1`，有效期设置为 5 分钟。到期前 1 纳秒请求 `/leases/session-1` 时会正常返回 `200` 和 `active`，但恰好在 12:05:00 请求时仍返回相同内容。

到期时刻应该视为已经过期，接口应返回 `404 Not Found`。到期前的读取、缺少键、不支持的请求方法、值拷贝和参数校验行为应保持不变。

## 复现命令

```powershell
go test . -run '^TestHandlerReturnsNotFoundAtExactExpiry$' -count=20
```

## 完整验证

```powershell
go test ./...
go vet ./...
go build ./...
```
