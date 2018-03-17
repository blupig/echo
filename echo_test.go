package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	testMode = true
	main()
}

func TestRootHandler(t *testing.T) {
	// Root path
	// Make request
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Call handler
	rootHandler(w, r)

	// Validate result
	res := w.Result()
	assert.Equal(t, 200, res.StatusCode, "Status code should be 200")

	// Should handle non-exist routes
	r = httptest.NewRequest("GET", "/test404", nil)
	w = httptest.NewRecorder()
	rootHandler(w, r)
	res = w.Result()
	assert.Equal(t, 404, res.StatusCode, "Status code should be 404")
}

func TestCacheHandler(t *testing.T) {
	// Make request
	r := httptest.NewRequest("GET", "/cache", nil)
	w := httptest.NewRecorder()

	// Call handler
	cacheHandler(w, r)

	// Validate result
	res := w.Result()
	assert.Equal(t, 200, res.StatusCode, "Status code should be 200")
	assert.Equal(t, "public, max-age=10", res.Header.Get("Cache-Control"), "Should have correct Cache-Control header")
}

func TestCPUHandler(t *testing.T) {
	// Make request
	r := httptest.NewRequest("GET", "/cpu", nil)
	w := httptest.NewRecorder()

	// Call handler
	cpuHandler(w, r)

	// Validate result
	res := w.Result()
	assert.Equal(t, 200, res.StatusCode, "Status code should be 200")
}

func TestExitHandler(t *testing.T) {
	// No X-Exit-Token
	r := httptest.NewRequest("GET", "/exit", nil)
	w := httptest.NewRecorder()
	exitHandler(w, r)
	res := w.Result()
	assert.Equal(t, 401, res.StatusCode, "Status code should be 401")

	// Empty token
	r = httptest.NewRequest("GET", "/exit", nil)
	r.Header.Set("X-Exit-Token", "")
	w = httptest.NewRecorder()
	exitHandler(w, r)
	res = w.Result()
	assert.Equal(t, 401, res.StatusCode, "Status code should be 401")

	// Wrong token
	exitToken = "real-token"
	r = httptest.NewRequest("GET", "/exit", nil)
	r.Header.Set("X-Exit-Token", "fake-token")
	w = httptest.NewRecorder()
	exitHandler(w, r)
	res = w.Result()
	assert.Equal(t, 401, res.StatusCode, "Status code should be 401")

	// Cannot test correct token as it exits the process
}

func TestHeadersHandler(t *testing.T) {
	// Make request
	r := httptest.NewRequest("GET", "/headers", nil)
	r.Header.Set("X-Test-Header", "test")
	w := httptest.NewRecorder()

	// Call handler
	headersHandler(w, r)

	// Validate result
	res := w.Result()
	body, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, 200, res.StatusCode, "Status code should be 200")
	assert.Equal(t, "{\n    \"X-Test-Header\": \"test\"\n}", string(body), "aaa")
}

func TestHealthHandler(t *testing.T) {
	// Make request
	r := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Call handler
	healthHandler(w, r)

	// Validate result
	res := w.Result()
	assert.Equal(t, 200, res.StatusCode, "Status code should be 200")
}

func TestIPHandler(t *testing.T) {
	// -- Without XFF header --
	// Make request
	r := httptest.NewRequest("GET", "/ip", nil)
	r.RemoteAddr = "10.0.1.2:1234"
	w := httptest.NewRecorder()

	// Call handler
	ipHandler(w, r)

	// Validate results
	res := w.Result()
	assert.Equal(t, 200, res.StatusCode, "Status code should be 200")

	var respObj map[string]interface{}
	body, _ := ioutil.ReadAll(res.Body)
	err := json.Unmarshal(body, &respObj)
	assert.NoError(t, err, "Should return valid JSON")
	assert.Equal(t, "10.0.1.2:1234", respObj["remote_addr"], "Should return correct IP")

	// -- With XFF header --
	// Make request
	r.Header.Set("X-Forwarded-For", "10.2.3.4:22222")
	w = httptest.NewRecorder()

	// Call handler
	ipHandler(w, r)

	// Validate results
	res = w.Result()
	assert.Equal(t, 200, res.StatusCode, "Status code should be 200")

	body, _ = ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &respObj)
	assert.NoError(t, err, "Should return valid JSON")
	assert.Equal(t, "10.2.3.4:22222", respObj["remote_addr"], "Should return correct IP")
}
