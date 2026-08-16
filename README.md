# appendledger

`appendledger` is a small Go service that keeps records contiguous within each stream.

Submit records with `POST /streams/{stream}/records` and a JSON body such as `{"sequence":1,"value":"created"}`. A record may be appended only at the next sequence, while repeating the same value is idempotent.

Run the checks with:

```text
go test ./...
go vet ./...
go build ./...
```
