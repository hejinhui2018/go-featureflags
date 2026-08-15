package rangeserve

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseRange(t *testing.T) {
	tests := []struct {
		name   string
		header string
		size   int64
		want   ByteRange
		err    error
	}{
		{name: "bounded", header: "bytes=2-5", size: 10, want: ByteRange{Start: 2, End: 5}},
		{name: "open ended", header: "bytes=7-", size: 10, want: ByteRange{Start: 7, End: 9}},
		{name: "suffix", header: "bytes=-4", size: 10, want: ByteRange{Start: 6, End: 9}},
		{name: "large suffix", header: "bytes=-20", size: 10, want: ByteRange{Start: 0, End: 9}},
		{name: "past end", header: "bytes=10-", size: 10, err: ErrUnsatisfiable},
		{name: "multiple", header: "bytes=0-1,4-5", size: 10, err: ErrInvalidRange},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRange(tt.header, tt.size)
			if !errors.Is(err, tt.err) {
				t.Fatalf("ParseRange() error = %v, want %v", err, tt.err)
			}
			if got != tt.want {
				t.Fatalf("ParseRange() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestHandlerServesSuffixRange(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/asset", nil)
	req.Header.Set("Range", "bytes=-4")
	recorder := httptest.NewRecorder()

	Handler([]byte("abcdefghij")).ServeHTTP(recorder, req)

	if got, want := recorder.Code, http.StatusPartialContent; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
	if got, want := recorder.Body.String(), "ghij"; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
	if got, want := recorder.Header().Get("Content-Range"), "bytes 6-9/10"; got != want {
		t.Fatalf("Content-Range = %q, want %q", got, want)
	}
}

func TestHandlerServesFullBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/asset", nil)
	recorder := httptest.NewRecorder()

	Handler([]byte("abcdefghij")).ServeHTTP(recorder, req)

	if got, want := recorder.Code, http.StatusOK; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
	if got, want := recorder.Body.String(), "abcdefghij"; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}

func TestHandlerRejectsUnsatisfiableRange(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/asset", nil)
	req.Header.Set("Range", "bytes=12-")
	recorder := httptest.NewRecorder()

	Handler([]byte("abcdefghij")).ServeHTTP(recorder, req)

	if got, want := recorder.Code, http.StatusRequestedRangeNotSatisfiable; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
	if got, want := recorder.Header().Get("Content-Range"), "bytes */10"; got != want {
		t.Fatalf("Content-Range = %q, want %q", got, want)
	}
}
