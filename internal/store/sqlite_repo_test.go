package store

import (
	"context"
	"database/sql"
	"testing"
	"time"
)

func openTestRepo(t *testing.T) Repository {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:?_journal_mode=WAL&_foreign_keys=on")
	if err != nil { t.Fatalf("open in-memory sqlite: %v", err) }
	db.SetMaxOpenConns(1)
	if err := migrate(db); err != nil { db.Close(); t.Fatalf("migrate: %v", err) }
	repo, err := NewSQLiteRepo(db)
	if err != nil { db.Close(); t.Fatalf("new sqlite repo: %v", err) }
	t.Cleanup(func() { repo.Close() })
	return repo
}

func TestCreateAndGetEndpoint(t *testing.T) {
	repo := openTestRepo(t)
	e, err := repo.CreateEndpoint(context.Background(), "test-one")
	if err != nil { t.Fatalf("CreateEndpoint: %v", err) }
	if e.ID == "" { t.Error("expected non-empty ID") }
	if e.URL != "test-one" { t.Errorf("url: got %q", e.URL) }
	if e.CreatedAt == "" || e.UpdatedAt == "" { t.Error("expected timestamps") }
	got, err := repo.GetEndpoint(context.Background(), e.ID)
	if err != nil { t.Fatalf("GetEndpoint: %v", err) }
	if got.ID != e.ID || got.URL != e.URL { t.Errorf("got %+v", got) }
}

func TestGetEndpointNotFound(t *testing.T) {
	repo := openTestRepo(t)
	_, err := repo.GetEndpoint(context.Background(), "nonexistent")
	if err == nil { t.Fatal("expected error") }
}

func TestCreateEndpointDuplicateURL(t *testing.T) {
	repo := openTestRepo(t)
	ctx := context.Background()
	if _, err := repo.CreateEndpoint(ctx, "dup"); err != nil { t.Fatal(err) }
	if _, err := repo.CreateEndpoint(ctx, "dup"); err == nil { t.Fatal("expected UNIQUE error") }
}

func TestAppendAndListRequests(t *testing.T) {
	repo := openTestRepo(t)
	ctx := context.Background()
	e, err := repo.CreateEndpoint(ctx, "reqs-test")
	if err != nil { t.Fatal(err) }
	r1 := &Request{Method: "POST", Path: "/webhook", Headers: `{"Content-Type":"application/json"}`, Query: "foo=bar", Body: []byte(`{"hello":"world"}`)}
	if err := repo.AppendRequest(ctx, e.ID, r1); err != nil { t.Fatal(err) }
	time.Sleep(time.Millisecond)
	r2 := &Request{Method: "GET", Path: "/webhook", Headers: `{}`}
	if err := repo.AppendRequest(ctx, e.ID, r2); err != nil { t.Fatal(err) }
	reqs, err := repo.ListRequests(ctx, e.ID, 10)
	if err != nil { t.Fatal(err) }
	if len(reqs) != 2 { t.Fatalf("got %d requests", len(reqs)) }
	if reqs[0].Method != "GET" { t.Errorf("first: got %s", reqs[0].Method) }
	if reqs[1].Method != "POST" { t.Errorf("second: got %s", reqs[1].Method) }
}

func TestListRequestsLimit(t *testing.T) {
	repo := openTestRepo(t)
	ctx := context.Background()
	e, _ := repo.CreateEndpoint(ctx, "limit-test")
	for i := 0; i < 5; i++ { repo.AppendRequest(ctx, e.ID, &Request{Method: "GET", Path: "/"}) }
	reqs, err := repo.ListRequests(ctx, e.ID, 3)
	if err != nil { t.Fatal(err) }
	if len(reqs) != 3 { t.Fatalf("got %d", len(reqs)) }
}

func TestListRequestsEmpty(t *testing.T) {
	repo := openTestRepo(t)
	ctx := context.Background()
	e, _ := repo.CreateEndpoint(ctx, "empty-test")
	reqs, _ := repo.ListRequests(ctx, e.ID, 10)
	if len(reqs) != 0 { t.Errorf("got %d", len(reqs)) }
}

func TestAppendRequestBodyBinary(t *testing.T) {
	repo := openTestRepo(t)
	ctx := context.Background()
	e, _ := repo.CreateEndpoint(ctx, "binary-test")
	body := []byte{0x00, 0x01, 0x02, 0xFF}
	if err := repo.AppendRequest(ctx, e.ID, &Request{Method: "POST", Path: "/binary", Body: body}); err != nil { t.Fatal(err) }
	reqs, _ := repo.ListRequests(ctx, e.ID, 1)
	if len(reqs) != 1 { t.Fatal("expected 1 request") }
	if len(reqs[0].Body) != 4 { t.Errorf("body len: %d", len(reqs[0].Body)) }
}

func TestMigrateIdempotent(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:?_journal_mode=WAL&_foreign_keys=on")
	if err != nil { t.Fatal(err) }
	defer db.Close()
	db.SetMaxOpenConns(1)
	if err := migrate(db); err != nil { t.Fatal(err) }
	if err := migrate(db); err != nil { t.Fatal(err) }
	var name string
	if err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='endpoints'").Scan(&name); err != nil {
		t.Fatalf("endpoints table: %v", err)
	}
	if err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='requests'").Scan(&name); err != nil {
		t.Fatalf("requests table: %v", err)
	}
}
