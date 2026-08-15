# Tenant Cache Scope Reproduction

## Behavior

Cached flag values can cross tenant boundaries when two tenants use the same
flag name. The first lookup populates the cache, and the second lookup can
receive that cached value even when its stored configuration is different.

## Reproduction

Run the focused regression test from the repository root:

```sh
go test ./internal/flags -run '^TestServiceEnabledSeparatesTenantCaches$' -count=1
```

## Expected result

`tenant-a/checkout` returns `true` and `tenant-b/checkout` returns `false`.

## Observed result

The base revision fails with:

```text
Enabled(tenant-b, checkout) = true, want false
```
