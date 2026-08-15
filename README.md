# checkpointstore

`checkpointstore` is a small Go service that records the latest processed offset for each stream.

Submit a checkpoint with `POST /checkpoints/{stream}` and a JSON body such as `{"offset":42}`. A new or advanced checkpoint returns `201 Created`, repeating the current offset returns `200 OK`, and moving backwards returns `409 Conflict`.

Run the checks with:

```text
go test ./...
go vet ./...
go build ./...
```
