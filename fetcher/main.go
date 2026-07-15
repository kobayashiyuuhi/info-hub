// info-hub fetcher: feeds.yml に定義された RSS/Atom フィードを並行取得し、
// マージ・日付降順ソート・件数制限した JSON を public/data/ に出力する。
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	feedsPath := flag.String("feeds", "feeds.yml", "path to feeds.yml")
	outDir := flag.String("out", "public/data", "output directory")
	limit := flag.Int("limit", 500, "max number of items to keep")
	flag.Parse()

	cfg, err := LoadConfig(*feedsPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	if len(cfg.Feeds) == 0 {
		log.Fatal("no feeds defined in config")
	}

	items, meta := FetchAll(cfg.Feeds)
	items = SortAndLimit(items, *limit)

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		log.Fatalf("failed to create output dir: %v", err)
	}
	if err := writeJSON(filepath.Join(*outDir, "items.json"), items); err != nil {
		log.Fatalf("failed to write items.json: %v", err)
	}
	if err := writeJSON(filepath.Join(*outDir, "meta.json"), meta); err != nil {
		log.Fatalf("failed to write meta.json: %v", err)
	}

	okCount := 0
	for _, f := range meta.Feeds {
		if f.Status == "ok" {
			okCount++
		}
	}
	fmt.Printf("done: %d items from %d/%d feeds -> %s\n", len(items), okCount, len(cfg.Feeds), *outDir)
}

func writeJSON(path string, v any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "")
	return enc.Encode(v)
}
