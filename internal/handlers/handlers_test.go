package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rizzza/echoserver/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	cfg = config.Get()
)

func TestRootRoute(t *testing.T) {
	t.Run("/", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		rr := httptest.NewRecorder()

		srv := New(cfg)
		srv.ServeHTTP(rr, req)

		resp := rr.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestEcho(t *testing.T) {
	t.Run("errors with no req content-type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/echo", http.NoBody)
		rr := httptest.NewRecorder()

		srv := New(cfg)
		srv.ServeHTTP(rr, req)

		resp := rr.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		bytes, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		require.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("content-type"))
		assert.Equal(t, strings.Contains(string(bytes), "invalid content-type"), true)
	})

	t.Run("errors with no json payload", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/echo", http.NoBody)
		req.Header.Set("content-type", "application/json")
		rr := httptest.NewRecorder()

		srv := New(cfg)
		srv.ServeHTTP(rr, req)

		resp := rr.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		bytes, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		require.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("content-type"))
		assert.Equal(t, strings.Contains(string(bytes), "invalid json"), true)
	})

	t.Run("errors with no 'echo' field", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(`{}`))
		req.Header.Set("content-type", "application/json")
		rr := httptest.NewRecorder()

		srv := New(cfg)
		srv.ServeHTTP(rr, req)

		resp := rr.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		bytes, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		require.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("content-type"))
		assert.Equal(t, strings.Contains(string(bytes), "field is required"), true)
	})

	t.Run("successful echo", func(t *testing.T) {
		payload := `{
			"echo":"hello there"
		}`

		req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(payload))
		req.Header.Set("content-type", "application/json")
		rr := httptest.NewRecorder()

		srv := New(cfg)
		srv.ServeHTTP(rr, req)

		resp := rr.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("content-type"))
		assert.Equal(t, `{"echo":"Echo server says -- hello there"}`, rr.Body.String())
	})
}
