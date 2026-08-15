# 超额批次留下部分预留

## 问题描述

容量为 10 时，向 `/reservations/batch` 提交包含 `alpha=4` 和 `beta=7` 的批次，接口返回 `409 Conflict`，但随后状态中已经有 4 个单位分配给了 alpha。

容量不足时，整个批次应该被拒绝，已用量保持不变；容量恰好用满的批次应该成功。请求校验、不同名称的隔离和现有状态码约定也应保持不变。

## 复现命令

```powershell
go test . -run '^TestHandlerRejectsBatchWithoutPartialReservations$' -count=20
```

## 完整验证

```powershell
go test ./...
go vet ./...
go build ./...
```
