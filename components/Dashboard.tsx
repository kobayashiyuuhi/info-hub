'use client';

import { useEffect, useMemo, useState } from 'react';
import type { FeedItem, Meta } from '@/lib/types';
import ItemCard from '@/components/ItemCard';

const BASE = process.env.NEXT_PUBLIC_BASE_PATH ?? '';
const READ_KEY = 'info-hub:read';
const FAV_KEY = 'info-hub:favs';
const MAX_READ_IDS = 2000;

function loadSet(key: string): Set<string> {
  if (typeof window === 'undefined') return new Set();
  try {
    const raw = localStorage.getItem(key);
    return raw ? new Set(JSON.parse(raw) as string[]) : new Set();
  } catch {
    return new Set();
  }
}

function saveSet(key: string, set: Set<string>, max = MAX_READ_IDS) {
  try {
    localStorage.setItem(key, JSON.stringify([...set].slice(-max)));
  } catch {
    /* quota exceeded etc. — ignore */
  }
}

export default function Dashboard() {
  const [items, setItems] = useState<FeedItem[]>([]);
  const [meta, setMeta] = useState<Meta | null>(null);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [category, setCategory] = useState<string>('all');
  const [query, setQuery] = useState('');
  const [unreadOnly, setUnreadOnly] = useState(false);
  const [favsOnly, setFavsOnly] = useState(false);
  const [readIds, setReadIds] = useState<Set<string>>(new Set());
  const [favIds, setFavIds] = useState<Set<string>>(new Set());

  useEffect(() => {
    setReadIds(loadSet(READ_KEY));
    setFavIds(loadSet(FAV_KEY));
    Promise.all([
      fetch(`${BASE}/data/items.json`).then((r) => {
        if (!r.ok) throw new Error(`items.json: HTTP ${r.status}`);
        return r.json();
      }),
      fetch(`${BASE}/data/meta.json`).then((r) => (r.ok ? r.json() : null)),
    ])
      .then(([itemsData, metaData]) => {
        setItems(itemsData as FeedItem[]);
        setMeta(metaData as Meta | null);
      })
      .catch((e) => setLoadError(String(e)));
  }, []);

  const categories = useMemo(() => {
    const set = new Set<string>();
    for (const it of items) set.add(it.category);
    return ['all', ...[...set].sort()];
  }, [items]);

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase();
    return items.filter((it) => {
      if (category !== 'all' && it.category !== category) return false;
      if (unreadOnly && readIds.has(it.id)) return false;
      if (favsOnly && !favIds.has(it.id)) return false;
      if (q) {
        const hay = `${it.title} ${it.source} ${it.summary ?? ''}`.toLowerCase();
        if (!hay.includes(q)) return false;
      }
      return true;
    });
  }, [items, category, query, unreadOnly, favsOnly, readIds, favIds]);

  const markRead = (id: string) => {
    setReadIds((prev) => {
      const next = new Set(prev);
      next.add(id);
      saveSet(READ_KEY, next);
      return next;
    });
  };

  const toggleFav = (id: string) => {
    setFavIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      saveSet(FAV_KEY, next);
      return next;
    });
  };

  const markAllRead = () => {
    setReadIds((prev) => {
      const next = new Set(prev);
      for (const it of filtered) next.add(it.id);
      saveSet(READ_KEY, next);
      return next;
    });
  };

  const unreadCount = items.filter((it) => !readIds.has(it.id)).length;
  const errorFeeds = meta?.feeds.filter((f) => f.status === 'error') ?? [];

  return (
    <div className="mx-auto max-w-3xl px-4 py-6">
      <header className="mb-6">
        <div className="flex items-baseline justify-between gap-2 flex-wrap">
          <h1 className="text-2xl font-bold tracking-tight">
            info-hub
            <span className="ml-3 text-sm font-normal text-slate-400">
              未読 {unreadCount} / 全 {items.length} 件
            </span>
          </h1>
          {meta && (
            <span className="text-xs text-slate-500">
              最終更新: {new Date(meta.generatedAt).toLocaleString('ja-JP')}
            </span>
          )}
        </div>
        {errorFeeds.length > 0 && (
          <p className="mt-2 rounded border border-amber-700/50 bg-amber-900/20 px-3 py-1.5 text-xs text-amber-400">
            取得失敗: {errorFeeds.map((f) => f.name).join(', ')}
          </p>
        )}
      </header>

      <div className="mb-4 flex flex-wrap items-center gap-2">
        {categories.map((c) => (
          <button
            key={c}
            onClick={() => setCategory(c)}
            className={`rounded-full px-3 py-1 text-sm transition-colors ${
              category === c
                ? 'bg-sky-600 text-white'
                : 'bg-slate-800 text-slate-300 hover:bg-slate-700'
            }`}
          >
            {c}
          </button>
        ))}
      </div>

      <div className="mb-6 flex flex-wrap items-center gap-3">
        <input
          type="search"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="キーワード検索..."
          className="min-w-48 flex-1 rounded-md border border-slate-700 bg-slate-900 px-3 py-1.5 text-sm placeholder-slate-500 focus:border-sky-600 focus:outline-none"
        />
        <label className="flex cursor-pointer items-center gap-1.5 text-sm text-slate-300">
          <input
            type="checkbox"
            checked={unreadOnly}
            onChange={(e) => setUnreadOnly(e.target.checked)}
            className="accent-sky-600"
          />
          未読のみ
        </label>
        <label className="flex cursor-pointer items-center gap-1.5 text-sm text-slate-300">
          <input
            type="checkbox"
            checked={favsOnly}
            onChange={(e) => setFavsOnly(e.target.checked)}
            className="accent-amber-500"
          />
          ★のみ
        </label>
        <button
          onClick={markAllRead}
          className="rounded-md border border-slate-700 px-2.5 py-1 text-xs text-slate-400 hover:bg-slate-800"
        >
          表示中を既読に
        </button>
      </div>

      {loadError && (
        <p className="rounded border border-red-800 bg-red-950/40 px-3 py-2 text-sm text-red-400">
          データ読み込みエラー: {loadError}
        </p>
      )}

      <ul className="space-y-2">
        {filtered.map((it) => (
          <ItemCard
            key={it.id}
            item={it}
            isRead={readIds.has(it.id)}
            isFav={favIds.has(it.id)}
            onRead={markRead}
            onToggleFav={toggleFav}
          />
        ))}
      </ul>

      {!loadError && filtered.length === 0 && items.length > 0 && (
        <p className="py-10 text-center text-sm text-slate-500">条件に一致する記事がありません</p>
      )}
      {!loadError && items.length === 0 && (
        <p className="py-10 text-center text-sm text-slate-500">読み込み中...</p>
      )}
    </div>
  );
}
