# 负 TTL 写入返回成功

## 现象

调用 `Put("new", "value", -time.Second)` 会返回成功，但记录随后立即读不到。已有 `theme=dark` 时调用 `Put("theme", "light", -time.Second)` 也会返回成功，原值被覆盖后立即过期。

## 复现

```text
go test ./ttlstore -run '^TestPutRejectsNegativeTTLWithoutChangingStore$' -count=20
```

修复前，负 TTL 没有返回 `ErrInvalidTTL`，并且会写入存储。

## 期望

- 任何负 TTL 都返回 `ErrInvalidTTL`。
- 被拒绝的请求不新增记录，也不改变已有记录。
- TTL 为 0 的记录仍不过期，正 TTL 仍按设定时间过期。
- 空 key 校验、正常更新和删除行为保持不变。

## 完整验证

```text
go test ./...
go vet ./...
go build ./...
```
