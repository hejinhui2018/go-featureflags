# 批量导入失败后保留原有数据

## 现象

存储原本包含 `theme=dark`。一次 JSONL 导入的第一行把它改成 `light`，第二行 JSON 格式错误。`Import` 返回第二行错误后，读取 `theme` 得到的却是 `light`，说明失败批次的前半部分已经生效。

## 复现

```text
go test ./batchstore -run '^Test(InvalidBatchLeavesStoreUnchanged|ReaderFailureLeavesStoreUnchanged)$' -count=20
```

修复前，第二行解析失败时已有值会被第一行覆盖。输入读取中途报错时也会留下此前解析成功的记录。

## 期望

- 任意一行 JSON 格式错误或 key 为空白时，原有记录不变，本批次的新记录不存在。
- 输入读取中途失败时同样不修改存储。
- 错误信息保留出错行号和底层错误。
- 全部合法的批量导入及已有值覆盖行为保持不变。

## 完整验证

```text
go test ./...
go vet ./...
go build ./...
```
