# 部分写入失败后的日志恢复

## 现象

日志已经成功保存 `theme=dark` 后，再写入 `layout=compact`。如果底层写入只完成一半并返回 `disk full`，`Put` 会正确返回错误，但关闭存储后再次调用 `Open` 会返回 `read journal: unexpected end of JSON input`。这会让先前已成功保存的数据也无法读取。

## 复现

在与 `journal` 包相同的包测试中，将 `Store.write` 替换为一个写入半条记录后返回错误的函数，然后执行一次 `Put`、关闭并重新打开日志。修复前，重新打开稳定失败。

```text
go test ./journal -run '^TestFailedAppendKeepsJournalReadable$' -count=20
```

## 期望

- 失败的 `Put` 返回底层写入错误。
- 重新打开日志成功，已有的 `theme=dark` 仍可读取。
- 写入失败的 `layout` 不存在。
- 正常写入、重启恢复和空白 key 校验保持原有行为。

## 完整验证

```text
go test ./...
go vet ./...
go build ./...
```
