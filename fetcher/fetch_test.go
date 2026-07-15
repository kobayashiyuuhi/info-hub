package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "feeds.yml")
	content := `feeds:
  - name: Example
    url: https://example.com/rss
    category: tech
  - name: Another
    url: https://example.org/atom.xml
    category: ai
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if len(cfg.Feeds) != 2 {
		t.Fatalf("expected 2 feeds, got %d", len(cfg.Feeds))
	}
	if cfg.Feeds[0].Name != "Example" || cfg.Feeds[0].Category != "tech" {
		t.Errorf("unexpected first feed: %+v", cfg.Feeds[0])
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	if _, err := LoadConfig(filepath.Join(t.TempDir(), "nope.yml")); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSortAndLimit(t *testing.T) {
	base := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	items := []Item{
		{ID: "a", Published: base.Add(1 * time.Hour)},
		{ID: "b", Published: base.Add(3 * time.Hour)},
		{ID: "a", Published: base.Add(1 * time.Hour)}, // duplicate
		{ID: "c", Published: base.Add(2 * time.Hour)},
	}
	got := SortAndLimit(items, 2)
	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}
	if got[0].ID != "b" || got[1].ID != "c" {
		t.Errorf("wrong order: %s, %s", got[0].ID, got[1].ID)
	}
}

func TestSortAndLimitNoLimit(t *testing.T) {
	items := []Item{{ID: "a"}, {ID: "b"}}
	if got := SortAndLimit(items, 0); len(got) != 2 {
		t.Errorf("limit 0 should keep all, got %d", len(got))
	}
}

func TestStripHTML(t *testing.T) {
	in := `<p>Hello <b>world</b></p> and <a href="x">link</a>`
	got := StripHTML(in)
	want := "Hello world and link"
	if got != want {
		t.Errorf("StripHTML = %q, want %q", got, want)
	}
}

func TestTruncate(t *testing.T) {
	if got := Truncate("こんにちは世界", 5); got != "こんにちは…" {
		t.Errorf("Truncate = %q", got)
	}
	if got := Truncate("short", 10); got != "short" {
		t.Errorf("Truncate = %q", got)
	}
}

func TestItemID(t *testing.T) {
	a := itemID("https://f", "https://l", "guid1", "t")
	b := itemID("https://f", "https://l", "guid1", "t")
	c := itemID("https://f", "https://l", "guid2", "t")
	if a != b {
		t.Error("same input should give same id")
	}
	if a == c {
		t.Error("different guid should give different id")
	}
	if len(a) != 16 {
		t.Errorf("id length = %d, want 16", len(a))
	}
}
