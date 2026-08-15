# 每 IP 限流复现

## 可观察行为

同一个客户端 IP 连续请求时，只要源端口发生变化，每次请求都会被当成新的客户端，从而绕过每 IP 一次的请求限制。不同客户端地址仍应拥有彼此独立的限额。

## 复现命令

```powershell
go test ./... -run '^TestSameIPv4AddressSharesLimitAcrossPorts$' -count=20
```

基线中第二个请求每次都返回 204。期望同一 IP 的第二个请求返回 429，无论源端口是否变化，同时 IPv4 和 IPv6 地址都遵循相同行为。
