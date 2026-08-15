# 同一个幂等键不应接受不同的请求体

## 问题描述

向 `/jobs` 连续发送两个请求时，如果它们使用相同的 `Idempotency-Key`，但请求体分别是 `alpha` 和 `beta`，第二个请求仍会返回第一次创建成功的 `201` 和 `job-1`。

期望第二个请求返回 `409 Conflict`，并且执行器只运行一次。相同幂等键和相同请求体的重复请求仍应正常回放第一次的结果；缺少幂等键的请求仍应返回 `400 Bad Request`。

## 复现命令

```powershell
go test . -run '^TestHandlerReturnsConflictForChangedPayload$' -count=20
```

## 完整验证

```powershell
go test ./...
go vet ./...
go build ./...
```
