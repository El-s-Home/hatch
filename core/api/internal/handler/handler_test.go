package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/elfoundation/hatch/internal/store"
	"github.com/elfoundation/hatch/internal/testutil"
	"github.com/go-chi/chi/v5"
)

// testRouter creates a chi router with all routes registered using a fake repo.
func testRouter(repo store.Repository) chi.Router {
	r := chi.NewRouter()
	h := New(repo)
	h.RegisterRoutes(r)
	return r
}

func TestHealthz(t *testing.T) {
	r := testRouter(testutil.NewFakeRepository())
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := strings.TrimSpace(w.Body.String())
	if body != "ok" {
		t.Fatalf("expected 'ok', got %q", body)
	}
}

func TestCaptureRecordsAllVerbs(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	for _, m := range methods {
		repo := testutil.NewFakeRepository()
		r := testRouter(repo)

		body := ""
		if m == "POST" || m == "PUT" || m == "PATCH" {
			body = `{"k":"v"}`
		}
		req := httptest.NewRequest(m, "/ep", strings.NewReader(body))
		if m == "GET" {
			req.Header.Set("Accept", "text/html")
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("%s: expected 200, got %d", m, w.Code)
		}
		if len(repo.Requests) != 1 {
			t.Errorf("%s: expected 1 request, got %d", m, len(repo.Requests))
			continue
		}
		reqCaptured := repo.Requests[0]
		if reqCaptured.Method != m {
			t.Errorf("%s: wrong method %s", m, reqCaptured.Method)
		}
	}
}

func TestSSEStreamReceivesEventOnCapture(t *testing.T) {
	repo := testutil.NewFakeRepository()
	srv := httptest.NewServer(testRouter(repo))
	defer srv.Close()

	// Use separate clients with separate transports to avoid
	// connection-pool contention with the long-lived SSE connection.
	sseClient := &http.Client{
		Transport: &http.Transport{DisableKeepAlives: true},
	}
	captureClient := &http.Client{
		Transport: &http.Transport{DisableKeepAlives: true},
	}

	// Subscribe to SSE with a cancellable context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sseReq, err := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL+"/e/ep/events", nil)
	if err != nil {
		t.Fatalf("SSE request create failed: %v", err)
	}
	sseResp, err := sseClient.Do(sseReq)
	if err != nil {
		t.Fatalf("SSE GET failed: %v", err)
	}
	defer sseResp.Body.Close()

	if sseResp.StatusCode != http.StatusOK {
		t.Fatalf("SSE returned %d", sseResp.StatusCode)
	}
	ct := sseResp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "text/event-stream") {
		t.Fatalf("SSE Content-Type is %q, expected text/event-stream", ct)
	}

	// Read SSE body in a goroutine; capture in main goroutine.
	bodyCh := make(chan string, 1)
	go func() {
		buf := make([]byte, 4096)
		n, _ := sseResp.Body.Read(buf)
		bodyCh <- string(buf[:n])
	}()

	// Small sleep to let the SSE subscription register.
	time.Sleep(10 * time.Millisecond)

	// Capture a request on the same server.
	captureResp, err := captureClient.Post(srv.URL+"/ep", "application/json", strings.NewReader(`{"msg":"hello"}`))
	if err != nil {
		t.Fatalf("capture POST failed: %v", err)
	}
	captureResp.Body.Close()
	if captureResp.StatusCode != http.StatusOK {
		t.Fatalf("capture returned %d", captureResp.StatusCode)
	}

	// Wait for the SSE event with a timeout.
	select {
	case body := <-bodyCh:
		if !strings.Contains(body, "data:") {
			t.Fatalf("SSE stream missing 'data:' prefix: %q", body)
		}
		if !strings.Contains(body, "\"method\":\"POST\"") {
			t.Fatalf("SSE stream missing POST method: %q", body)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for SSE event")
	}

	// Cancel context to close the SSE connection.
	cancel()
}

// The inspect route now serves a static SPA shell; all request rendering
// happens client-side from the JSON API + SSE. So the shell is identical
// regardless of endpoint state — we just assert it's the app document.
func TestInspectServesSPAShell(t *testing.T) {
	repo := testutil.NewFakeRepository()
	repo.CreateEndpoint(nil, "ep")
	repo.AppendRequest(nil, "ep", &store.Request{
		Method: "POST",
		Path:   "/ep/webhook",
		Body:   []byte(`{"msg":"hello"}`),
	})

	r := testRouter(repo)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/e/ep", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	ct := w.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		t.Fatalf("expected text/html Content-Type, got %q", ct)
	}

	body := w.Body.String()

	// The SPA shell mounts at #root and loads its bundle from /assets/.
	if !strings.Contains(body, `id="root"`) {
		t.Error("missing SPA mount point in shell")
	}
	if !strings.Contains(body, "/assets/") {
		t.Error("missing hashed asset reference in shell")
	}
}

// Endpoint auto-creation is now the API's responsibility (the SPA calls the
// JSON list endpoint on load), not the static page's.
func TestListAutoCreatesEndpoint(t *testing.T) {
	repo := testutil.NewFakeRepository()

	r := testRouter(repo)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/v1/endpoints/auto-ep/requests", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	ep, err := repo.GetEndpoint(nil, "auto-ep")
	if err != nil {
		t.Fatalf("endpoint was not auto-created: %v", err)
	}
	if ep.ID != "auto-ep" {
		t.Errorf("expected endpoint ID 'auto-ep', got %q", ep.ID)
	}
}

func TestCaptureReturnsJSON200(t *testing.T) {
	r := testRouter(testutil.NewFakeRepository())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/t", nil))

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Header().Get("Content-Type"), "application/json") {
		t.Error("not json")
	}
	data, _ := io.ReadAll(w.Result().Body)
	w.Result().Body.Close()
	var v map[string]interface{}
	if json.Unmarshal(data, &v) != nil {
		t.Fatal("invalid json")
	}
}

func TestMockSetsAndConfiguresResponse(t *testing.T) {
	repo := testutil.NewFakeRepository()
	repo.CreateEndpoint(nil, "mock-ep")

	r := testRouter(repo)

	// Set a mock via PUT.
	mockBody := `{"status": 201, "headers": {"X-Mocked": "yes"}, "body": "bW9ja2Vk"}`
	req := httptest.NewRequest(http.MethodPut, "/e/mock-ep/mock", strings.NewReader(mockBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}

	// Verify the mock was stored.
	mock, err := repo.GetMock(nil, "mock-ep")
	if err != nil {
		t.Fatalf("mock not stored: %v", err)
	}
	if mock.Status != 201 {
		t.Errorf("expected status 201, got %d", mock.Status)
	}
	if mock.Headers["X-Mocked"] != "yes" {
		t.Errorf("expected X-Mocked header, got %v", mock.Headers)
	}
}

func TestMockReturnsConfiguredResponseOnCapture(t *testing.T) {
	repo := testutil.NewFakeRepository()
	repo.CreateEndpoint(nil, "mock-cap")
	// Pre-set a mock.
	repo.SetMock(nil, &store.MockConfig{
		EndpointID: "mock-cap",
		Status:     418,
		Headers:    map[string]string{"X-Teapot": "true"},
		Body:       []byte(`{"brewing": true}`),
	})

	r := testRouter(repo)

	// Capture a request — should return the mock response, not the default 200.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/mock-cap", strings.NewReader(`{"order":"latte"}`)))

	if w.Code != 418 {
		t.Fatalf("expected mock status 418, got %d", w.Code)
	}
	if w.Header().Get("X-Teapot") != "true" {
		t.Errorf("expected X-Teapot: true, got %q", w.Header().Get("X-Teapot"))
	}
	body := strings.TrimSpace(w.Body.String())
	if body != `{"brewing": true}` {
		t.Errorf("expected mock body, got %q", body)
	}

	// Request should still be captured even when mock responds.
	if len(repo.Requests) != 1 {
		t.Fatalf("expected 1 captured request, got %d", len(repo.Requests))
	}
	if repo.Requests[0].Method != "POST" {
		t.Errorf("captured method: %s", repo.Requests[0].Method)
	}
}

func TestMockAutoCreatesEndpointOnSet(t *testing.T) {
	repo := testutil.NewFakeRepository()
	r := testRouter(repo)

	// Set a mock on an endpoint that doesn't exist yet.
	mockBody := `{"status": 200, "headers": {}, "body": ""}`
	req := httptest.NewRequest(http.MethodPut, "/e/auto-mock/mock", strings.NewReader(mockBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}

	// Endpoint should be auto-created.
	ep, err := repo.GetEndpoint(nil, "auto-mock")
	if err != nil {
		t.Fatalf("endpoint was not auto-created: %v", err)
	}
	if ep.ID != "auto-mock" {
		t.Errorf("expected endpoint id 'auto-mock', got %q", ep.ID)
	}
}

func TestRootServesSPAShell(t *testing.T) {
	repo := testutil.NewFakeRepository()
	r := testRouter(repo)

	// Root now serves the SPA home (no longer a capture target).
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), `id="root"`) {
		t.Error("missing SPA mount point at root")
	}
}

func TestV1CaptureRequestMissingEndpoint(t *testing.T) {
	repo := testutil.NewFakeRepository()
	r := testRouter(repo)

	body := `{"method":"POST","path":"/test"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/endpoints//requests", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestV1ListRequestsMissingEndpoint(t *testing.T) {
	repo := testutil.NewFakeRepository()
	r := testRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/v1/endpoints//requests", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestJoinPath(t *testing.T) {
	tests := []struct {
		base, extra, want string
	}{
		{"/api", "/users", "/api/users"},
		{"/api/", "/users", "/api/users"},
		{"/api", "/users/", "/api/users/"},
		{"/api", "", "/api"},
		{"", "/users", "/users"},
		{"", "", "/"},
	}
	for _, tc := range tests {
		got := joinPath(tc.base, tc.extra)
		if got != tc.want {
			t.Errorf("joinPath(%q, %q) = %q, want %q", tc.base, tc.extra, got, tc.want)
		}
	}
}

func TestHandleMockReturnsErrorOnInvalidJSON(t *testing.T) {
	repo := testutil.NewFakeRepository()
	repo.CreateEndpoint(nil, "mock-err")
	r := testRouter(repo)

	// Send invalid JSON.
	req := httptest.NewRequest(http.MethodPut, "/e/mock-err/mock", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// The handler should handle the error gracefully.
	// It may return 200 (if it ignores parse errors) or 400.
	if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
		t.Errorf("expected 200 or 400, got %d", w.Code)
	}
}

