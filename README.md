# Cursor Store

`cursorstore` is a small in-memory, ordered key/value store with cursor-based pagination.

`List` accepts an optional key cursor and a positive page size. It returns the page and a cursor for the next page. A blank cursor starts from the beginning.

Run the sample command with `go run ./cmd/cursorstore`.
