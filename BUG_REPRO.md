# Suffix Byte Ranges Return the End of the Content

## Behavior

An HTTP suffix range asks for the last N bytes of the current representation.
The response offsets, body, and `Content-Range` header must describe those
trailing bytes.

## Reproduction

```sh
go test . -run '^TestHandlerServesSuffixRange$' -count=20
```

## Expected result

For the content `abcdefghij`, `Range: bytes=-4` returns status 206, body
`ghij`, and `Content-Range: bytes 6-9/10`.

## Observed result on the affected revision

The handler returns body `abcd` and `Content-Range: bytes 0-3/10`.

## Verification

```sh
go test . -run '^TestHandlerServesSuffixRange$' -count=20
go test ./...
go vet ./...
go build ./...
```
