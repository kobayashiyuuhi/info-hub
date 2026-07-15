package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
	"gopkg.in/yaml.v3"
)

// FeedDef は feeds.yml の 1 エントリ。
type FeedDef struct {
	Name     string `yaml:"name" json:"name"`
	URL      string `yaml:"url" json:"url"`
	Category string `yaml:"category" json:"category"`
}

// Config は feeds.yml 全体。
type Config struct {
	Feeds []FeedDef `yaml:"feeds"`
}

// Item はフロントに渡す 1 記事。
type Item struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Source    string    `json:"source"`
	Category  string    `json:"category"`
	Published time.Time `json:"published"`
	Summary   string    `json:"summary,omitempty"`
}

// FeedStatus は各フィードの取得結果。
type FeedStatus struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	Category  string `json:"category"`
	Status    string `json:"status"` // "ok" | "error"
	Error     string `json:"error,omitempty"`
	ItemCount int    `json:"itemCount"`
}

// Meta は meta.json の内容。
type Meta struct {
	GeneratedAt time.Time    `json:"generatedAt"`
	Feeds       []FeedStatus `json:"feeds"`
}

const (
	fetchTimeout   = 20 * time.Second
	maxConcurrency = 8
	maxSummaryLen  = 300
)

// LoadConfig は feeds.yml を読み込む。
func LoadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// FetchAll は全フィードを並行取得する。失敗したフィードはスキップし meta に記録する。
func FetchAll(feeds []FeedDef) ([]Item, Meta) {
	var (
		mu       sync.Mutex
		items    []Item
		statuses = make([]FeedStatus, len(feeds))
		wg       sync.WaitGroup
		sem      = make(chan struct{}, maxConcurrency)
	)

	for i, fd := range feeds {
		wg.Add(1)
		go func(i int, fd FeedDef) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			st := FeedStatus{Name: fd.Name, URL: fd.URL, Category: fd.Category}
			got, err := fetchOne(fd)
			if err != nil {
				st.Status = "error"
				st.Error = err.Error()
			} else {
				st.Status = "ok"
				st.ItemCount = len(got)
				mu.Lock()
				items = append(items, got...)
				mu.Unlock()
			}
			statuses[i] = st
		}(i, fd)
	}
	wg.Wait()

	return items, Meta{GeneratedAt: time.Now().UTC(), Feeds: statuses}
}

func fetchOne(fd FeedDef) ([]Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), fetchTimeout)
	defer cancel()

	fp := gofeed.NewParser()
	fp.Client = &http.Client{Timeout: fetchTimeout}
	fp.UserAgent = "info-hub-fetcher/1.0 (+https://github.com/kobayashiyuuhi/info-hub)"

	feed, err := fp.ParseURLWithContext(fd.URL, ctx)
	if err != nil {
		return nil, err
	}

	items := make([]Item, 0, len(feed.Items))
	for _, it := range feed.Items {
		items = append(items, toItem(it, fd))
	}
	return items, nil
}

func toItem(it *gofeed.Item, fd FeedDef) Item {
	pub := time.Time{}
	if it.PublishedParsed != nil {
		pub = it.PublishedParsed.UTC()
	} else if it.UpdatedParsed != nil {
		pub = it.UpdatedParsed.UTC()
	}
	return Item{
		ID:        itemID(fd.URL, it.Link, it.GUID, it.Title),
		Title:     strings.TrimSpace(it.Title),
		URL:       strings.TrimSpace(it.Link),
		Source:    fd.Name,
		Category:  fd.Category,
		Published: pub,
		Summary:   Truncate(StripHTML(it.Description), maxSummaryLen),
	}
}

func itemID(feedURL, link, guid, title string) string {
	key := guid
	if key == "" {
		key = link
	}
	if key == "" {
		key = title
	}
	h := sha1.Sum([]byte(feedURL + "|" + key))
	return hex.EncodeToString(h[:])[:16]
}

// SortAndLimit は日付降順にソートし直近 limit 件に絞る。ID重複は除去する。
func SortAndLimit(items []Item, limit int) []Item {
	seen := make(map[string]bool, len(items))
	uniq := make([]Item, 0, len(items))
	for _, it := range items {
		if seen[it.ID] {
			continue
		}
		seen[it.ID] = true
		uniq = append(uniq, it)
	}
	sort.SliceStable(uniq, func(a, b int) bool {
		return uniq[a].Published.After(uniq[b].Published)
	})
	if limit > 0 && len(uniq) > limit {
		uniq = uniq[:limit]
	}
	return uniq
}

// StripHTML は雑にHTMLタグを除去する（要約表示用）。
func StripHTML(s string) string {
	var b strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
			b.WriteRune(' ')
		case !inTag:
			b.WriteRune(r)
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

// Truncate は rune 単位で最大 n 文字に切り詰める。
func Truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}
