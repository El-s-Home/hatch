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
)

type fakeRepo struct {
	endpoints map[string]*store.Endpoint
	requests  []*store.Request
}
func newFakeRepo() *fakeRepo { return &fakeRepo{endpoints: map[string]*store.Endpoint{}} }
func (f *fakeRepo) CreateEndpoint(_ context.Context, u string) (*store.Endpoint, error) {
	e := &store.Endpoint{ID: u, URL: u, CreatedAt: "t", UpdatedAt: "t"}
	f.endpoints[u] = e; return e, nil
}
func (f *fakeRepo) GetEndpoint(_ context.Context, id string) (*store.Endpoint, error) {
	e, ok := f.endpoints[id]; if !ok { return nil, errNotFound }; return e, nil
}
func (f *fakeRepo) AppendRequest(_ context.Context, eid string, r *store.Request) error {
	r.ID = "x"; r.EndpointID = eid; f.requests = append(f.requests, r); return nil
}
func (f *fakeRepo) ListRequests(_ context.Context, _ string, _ int) ([]*store.Request, error) { return f.requests, nil }
func (f *fakeRepo) Close() error { return nil }
var errNotFound = &se{"nf"}
type se struct{ m string }
func (e *se) Error() string { return e.m }

func TestRecordsAllVerbs(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	for _, m := range methods {
		repo := newFakeRepo()
		h := New(repo)
		body := ""
		if m == "POST" || m == "PUT" || m == "PATCH" { body = `{"k":"v"}` }
		req := httptest.NewRequest(m, "/ep", strings.NewReader(body))
		if m == "GET" { req.Header.Set("Accept", "text/html") }
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Result().StatusCode != 200 { t.Errorf("%s: expected 200", m) }
		if len(repo.requests) != 1 { t.Errorf("%s: expected 1 request", m); continue }
		r := repo.requests[0]
		if r.Method != m { t.Errorf("%s: wrong method %s", m, r.Method) }
		if r.Path != "/ep" { t.Errorf("%s: wrong path", m) }
	}
}

func TestHandlerReturnsJSON200(t *testing.T) {
	h := New(newFakeRepo())
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/t", nil))
	resp := w.Result()
	if resp.StatusCode != 200 { t.Fatalf("expected 200") }
	if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") { t.Error("not json") }
	data, _ := io.ReadAll(resp.Body); resp.Body.Close()
	var v map[string]interface{}
	if json.Unmarshal(data, &v) != nil { t.Fatal("invalid json") }
}

func TestExtractID(t *testing.T) {
	for _, tc := range []struct{ p, w string }{
		{"/a","a"},{"/a/b","a"},{"/",""},{"/x?y=z","x"},
	} {
		if g := extractID(tc.p); g != tc.w {
			t.Errorf("extractID(%q)=%q want %q", tc.p, g, tc.w)
		}
	}
}
