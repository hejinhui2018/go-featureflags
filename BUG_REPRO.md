# 空白 key 的错误语义

## 复现

启动服务后请求一个只包含空白字符的记录 key：

```text
GET /v1/records/%20%20
```

修复前服务返回 `404`，调用方无法区分无效 key 和合法但不存在的 key。

## 期望行为

只包含空白字符的 key 应返回 `400`，响应为：

```json
{"error":"invalid key"}
```

合法但不存在的 key 仍返回 `404` 和 `{"error":"not found"}`；已存在的 key 仍返回 `200` 及原有记录内容。

## 验证

```text
go test ./internal/httpapi -run '^TestGetRecordWhitespaceKey$' -count=20
go test ./...
go vet ./...
go build ./...
```
