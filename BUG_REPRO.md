# 满容量时更新已有记录

## 现象

容量为2的存储已经包含 `theme=dark` 和 `density=comfortable`。调用 `Put("theme", "light")` 返回 `store capacity reached`，已有记录无法更新。

## 复现

```text
go test ./limitstore -run '^Test(UpdateAtCapacity|UpdateAfterDeleteAndRefill|ZeroCapacityRejectsNewKey)$' -count=20
```

修复前，存储达到上限后，更新已有 key 会被当成新增 key 拒绝。

## 期望

- 满容量时更新已有 key 成功，记录数量不变。
- 满容量时写入新 key 仍返回 `ErrCapacity`，且不创建记录。
- 删除记录后可以新增，重新满容量后仍可更新已有 key。
- 空白 key 校验、正常读写和并发安全保持不变。

## 完整验证

```text
go test ./...
go vet ./...
go build ./...
```
