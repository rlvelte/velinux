package http

import (
	"io"
	nethttp "net/http"
)

const ContextKey = "http"

// HTTP is the default backend using the standard library.
type HTTP struct{}

// Get performs an GET request and returns the response body bytes.
func (h *HTTP) Get(url string) ([]byte, error) {
	resp, err := nethttp.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
