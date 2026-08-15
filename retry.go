package httpretrybody

import (
	"bytes"
	"io"
	"net/http"
)

// DoWithRetry retries a request once when the first response is temporary.
func DoWithRetry(client *http.Client, req *http.Request) (*http.Response, error) {
	if client == nil {
		client = http.DefaultClient
	}

	// Buffer the request body so it can be resent on retry.
	var bodyBytes []byte
	hasBody := req.Body != nil
	if hasBody {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			return nil, err
		}
	}

	for attempt := 0; attempt < 2; attempt++ {
		// Reset the body for each attempt so the payload is available on retry.
		if hasBody {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			req.GetBody = func() (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewReader(bodyBytes)), nil
			}
			req.ContentLength = int64(len(bodyBytes))
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode < 500 || attempt == 1 {
			return resp, nil
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
	return nil, io.ErrUnexpectedEOF
}
