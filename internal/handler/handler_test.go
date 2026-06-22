package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/elfoundation/hatch/internal/store"
	"github.com/go-chi/chi/v5"
)

// fakeRepo implements store.Repository for tests.
type fakeRepo struct {
	endpoints map[string]*store.Endpoint
	requests  []*store.Request
	mocks     map[string]*store.MockConfig
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		endpoints: map[string]*store.Endpoint{},
		mocks:     map[string]*store.MockConfig{},
	}
}

func (f *fakeRepo) CreateEndpoint(_ context.Context, u string) (*store.Endpoint, error) {
	e := &store.Endpoint{ID: u, URL: u, CreatedAt: "t", UpdatedAt: "t"}
	f.endpoints[u] = e
	return e, nil
}
func (f *fakeRepo) GetEndpoint(_ context.Context, id string) (*store.Endpoint, error) {
	e, ok := f.endpoints[id]
	if !ok {
		return nil, errNotFound
	}
	return e, nil
}
func (f *fakeRepo) AppendRequest(_ context.Context, eid string, r *store.Request) error {
	r.ID = "req-" + string(rune(len(f.requests)+'0'))
	r.EndpointID = eid
	f.requests = append(f.requests, r)
	return nil
}
func (f *fakeRepo) GetRequest(_ context.Context, id string) (*store.Request, error) {
	for _, r := range f.requests {
		if r.ID == id {
			return r, nil
		}
	}
	return nil, errNotFound
}
func (f *fakeRepo) ListRequests(_ context.Context, _ string, _ int) ([]*store.Request, error) {
	return f.requests, nil
}
func (f *fakeRepo) GetMock(_ context.Context, endpointID string) (*store.MockConfig, error) {
	m, ok := f.mocks[endpointID]
	if !ok {
		return nil, errNotFound
	}
	return m, nil
}
func (f *fakeRepo) SetMock(_ context.Context, mock *store.MockConfig) error {
	f.mocks[mock.EndpointID] = mock
	return nil
}
func (f *fakeRepo) Close() error { return nil }

var errNotFound = &se{"nf"}

type se struct{ m string }

func (e *se) Error() string { return e.m }

// testRouter creates a chi router with all routes registered using a fake repo.
func testRouter(repo store.Repository) chi.Router {
	r := chi.NewRouter()
	h := New(repo)
	h.RegisterRoutes(r)
	return r
}

func TestHealthz(t *testing.T) {
	r := testRouter(newFakeRepo())
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
		repo := newFakeRepo()
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
		if len(repo.requests) != 1 {
			t.Errorf("%s: expected 1 request, got %d", m, len(repo.requests))
			continue
		}
		reqCaptured := repo.requests[0]
		if reqCaptured.Method != m {
			t.Errorf("%s: wrong method %s", m, reqCaptured.Method)
		}
	}
}

func TestCaptureReturnsJSON200(t *testing.T) {
	r := testRouter(newFakeRepo())
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
