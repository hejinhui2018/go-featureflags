# 仅创建写入覆盖已有记录

## 现象

存储中已有 `theme=dark`、修订号为 1。调用 `Put("theme", "light", 0)` 仍然返回成功，并把记录改成 `theme=light`、修订号 2。调用方原本只想在 key 不存在时创建记录，却覆盖了已有数据。

## 复现

```text
go test ./revisionstore -run '^TestCreateOnlyRevisionConflictsWithExistingRecord$' -count=20
```

修复前，已有 key 会错误接受修订号 0 并完成更新。

## 期望

- 已有 key 收到 `expectedRevision=0` 时返回 `ErrConflict`。
- 冲突不改变原值或原修订号。
- 新 key 仍可用修订号 0 创建。
- 匹配修订号的更新、过期修订号冲突、缺失记录和输入校验保持不变。

## 完整验证

```text
go test ./...
go vet ./...
go build ./...
```
