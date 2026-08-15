# quotareservoir

`quotareservoir` is a small Go service that applies named reservations against a fixed capacity.

Submit a JSON object to `POST /reservations/batch`, for example `{"reservations":[{"name":"alpha","units":4}]}`. A batch is accepted only when every item can be applied without exceeding the remaining capacity.

Run the checks with:

```text
go test ./...
go vet ./...
go build ./...
```
