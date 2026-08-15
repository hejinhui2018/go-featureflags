# 重复提交检查点返回冲突

## 问题描述

先向 `/checkpoints/orders` 提交 `{"offset":42}` 时接口返回 `201 Created`，再次提交完全相同的 offset 时却返回 `409 Conflict`。

重复提交当前检查点应该是幂等操作并返回 `200 OK`。提交更小的 offset 应继续返回 `409 Conflict`，提交更大的 offset 应成功前进，现有请求校验和不同 stream 的隔离行为也应保持不变。

## 复现命令

```powershell
go test . -run '^TestHandlerAcceptsRepeatedCheckpoint$' -count=20
```

## 完整验证

```powershell
go test ./...
go vet ./...
go build ./...
```
