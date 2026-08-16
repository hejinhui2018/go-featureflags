# Request cancellation reproduction

Run the focused test from this directory:

```text
go test ./... -run '^TestMarkRequestPreservesCancellation$' -count=1
```

The test cancels a request before it enters the middleware. The downstream handler must observe `context.Canceled` immediately while still receiving the middleware marker. The second test confirms that a live request keeps its normal status and marker.
