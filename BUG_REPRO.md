# 分页 offset 超出末尾时服务崩溃

## 现象

存储中已有 `alpha` 和 `beta` 两条记录。调用 `List(5, 2)` 时进程触发 slice bounds panic，而不是把当前没有可返回记录当作空页处理。

## 复现

```text
go test ./pagestore -run '^TestListOffsetAfterLastRecordReturnsEmpty$' -count=20
```

修复前，每次调用都会稳定触发 panic。

## 期望

- offset 大于或等于当前记录数时返回非 nil 空页和 nil error。
- 正常分页和尾页截断保持不变。
- 负 offset 仍返回 `ErrInvalidOffset`，非正 limit 仍返回 `ErrInvalidLimit`。
- 越界分页读取不改变已有记录。

## 完整验证

```text
go test ./...
go vet ./...
go build ./...
```
