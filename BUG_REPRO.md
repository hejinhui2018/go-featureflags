# Request retry body reproduction

## Observable behavior

When a POST receives a temporary 5xx response, the retry returns success but the second request arrives without the original body. A successful request that does not need a retry is sent once.

## Reproduce

```powershell
go test ./... -run '^TestRetryResendsPostBody$' -count=20
```

On the baseline, the test fails because the retried request body is empty. The expected behavior is that the second attempt receives the same JSON payload as the first attempt.
