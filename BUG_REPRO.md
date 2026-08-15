# Canceled request limiter reproduction

## Observable behavior

After a canceled request returns from a concurrency-limited HTTP handler, the next normal request remains blocked even though no handler is running. Completed requests do not cause the same behavior.

## Reproduce

```powershell
go test ./... -run '^TestCanceledRequestReleasesSlot$' -count=20
```

On the baseline, every run times out waiting for the normal request. The expected behavior is that the canceled request frees its concurrency slot before returning so the next request can run.
