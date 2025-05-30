package api

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

func LogAndResetBody(t *testing.T, res *http.Response) {
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
		return
	}
	t.Logf("Response body: %s", string(bodyBytes))
	res.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
}
