# 未知游标返回第一页

## 现象

存储中包含 `alpha`、`beta` 和 `gamma` 三条记录。调用 `List("deleted-item", 2)` 时，返回了 `alpha`、`beta` 和下一页游标 `beta`，错误值为空。调用方会把第一页当作新一页再次处理。

## 复现

```text
go test ./cursorstore -run '^TestListRejectsUnknownCursor$' -count=20
```

修复前，该测试会看到第一页数据而不是 `ErrInvalidCursor`。

## 期望

- 非空游标在当前存储中不存在时返回 `ErrInvalidCursor`。
- 错误响应不包含条目或下一页游标。
- 空游标仍读取第一页，已有游标继续从对应条目之后读取。
- 页大小校验、写入更新和键排序行为保持不变。

## 完整验证

```text
go test ./...
go vet ./...
go build ./...
```
