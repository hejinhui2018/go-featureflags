# Batch Entries Use the Cumulative Balance

## Behavior

`ApplyBatch` evaluates entries in input order and commits all balance changes
only when every entry is valid.

## Reproduction

```sh
go test . -run '^TestApplyBatchIsAtomicOnError$' -count=20
```

## Expected result

When an account starts at 10 and a batch debits 6 followed by 5, the method
returns `ErrInsufficientFunds` and leaves the ledger unchanged.

## Observed result on the affected revision

The method returns `nil` because each debit is checked separately against the
starting balance. It then commits a balance of -1.

## Verification

```sh
go test . -run '^TestApplyBatchIsAtomicOnError$' -count=20
go test ./...
go vet ./...
go build ./...
```
