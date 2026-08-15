// Package rangeserve serves immutable byte content with HTTP Range support.
package rangeserve

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var (
	// ErrInvalidRange reports a malformed or unsupported Range header.
	ErrInvalidRange = errors.New("invalid range")
	// ErrUnsatisfiable reports a range outside the current representation.
	ErrUnsatisfiable = errors.New("unsatisfiable range")
)

// ByteRange is an inclusive range of byte offsets.
type ByteRange struct {
	Start int64
	End   int64
}

// ParseRange parses one HTTP bytes range for a representation of size bytes.
func ParseRange(header string, size int64) (ByteRange, error) {
	if size < 0 || !strings.HasPrefix(header, "bytes=") {
		return ByteRange{}, ErrInvalidRange
	}
	spec := strings.TrimSpace(strings.TrimPrefix(header, "bytes="))
	if spec == "" || strings.Contains(spec, ",") {
		return ByteRange{}, ErrInvalidRange
	}

	parts := strings.SplitN(spec, "-", 2)
	if len(parts) != 2 {
		return ByteRange{}, ErrInvalidRange
	}

	if parts[0] == "" {
		suffix, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil || suffix <= 0 {
			return ByteRange{}, ErrInvalidRange
		}
		if size == 0 {
			return ByteRange{}, ErrUnsatisfiable
		}
		if suffix > size {
			suffix = size
		}
		return ByteRange{Start: 0, End: suffix - 1}, nil
	}

	start, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || start < 0 {
		return ByteRange{}, ErrInvalidRange
	}
	if start >= size {
		return ByteRange{}, ErrUnsatisfiable
	}
	if parts[1] == "" {
		return ByteRange{Start: start, End: size - 1}, nil
	}

	end, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || end < start {
		return ByteRange{}, ErrInvalidRange
	}
	if end >= size {
		end = size - 1
	}
	return ByteRange{Start: start, End: end}, nil
}

// Handler returns an HTTP handler for content with single-range support.
func Handler(content []byte) http.Handler {
	stored := append([]byte(nil), content...)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Accept-Ranges", "bytes")
		rangeHeader := r.Header.Get("Range")
		if rangeHeader == "" {
			w.Header().Set("Content-Length", strconv.Itoa(len(stored)))
			_, _ = w.Write(stored)
			return
		}

		selected, err := ParseRange(rangeHeader, int64(len(stored)))
		if err != nil {
			w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", len(stored)))
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}

		body := stored[selected.Start : selected.End+1]
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", selected.Start, selected.End, len(stored)))
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		w.WriteHeader(http.StatusPartialContent)
		_, _ = w.Write(body)
	})
}
