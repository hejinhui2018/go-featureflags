# Merge Leaves Its Input Unchanged

## Behavior

`Merge` returns normalized, sorted closed intervals. Callers may reuse the
slice they pass in after the result is calculated.

## Reproduction

```sh
go test ./intervals -run '^TestMergeDoesNotModifyInput$' -count=20
```

## Expected result

The result is normalized, while the input slice keeps its original ordering
and endpoint values.

## Observed result on the base revision

When the input includes an out-of-order interval with reversed endpoints,
`Merge` returns the normalized result but also changes the input slice.

## Verification

```sh
go test ./intervals -run '^TestMergeDoesNotModifyInput$' -count=20
go test ./...
go vet ./...
go build ./...
```
