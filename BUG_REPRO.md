# 相邻闭区间合并

## 现象

调用 `intervals.Merge` 处理 `[1,3]` 和 `[3,5]` 时，结果仍然包含两段区间。两个闭区间在端点处相接时，调用方会把它们当成一个连续范围使用，因此不应保留中间边界。

## 期望

输入 `[1,3]`、`[3,5]` 应返回单段 `[1,5]`。已有的倒序端点归一化、排序和重叠区间合并行为保持不变。

## 复现

在基线提交 `b230e1ea1b24bb5dea0e9f5c503f7db1466780dc` 上运行：

```text
go test ./intervals -run '^TestMergeMergesTouchingIntervals$' -count=20
```

该命令稳定失败，实际结果为两段 `[1,3]`、`[3,5]`。

## 验证

修复后运行以下命令均通过：

```text
go test ./intervals -run '^TestMergeMergesTouchingIntervals$' -count=20
go test ./...
go vet ./...
go build ./...
```
