'use client';

import type { FeedItem } from '@/lib/types';

function timeAgo(iso: string): string {
  const t = new Date(iso).getTime();
  if (Number.isNaN(t) || t === 0) return '';
  const diff = Date.now() - t;
  const min = Math.floor(diff / 60000);
  if (min < 1) return 'たった今';
  if (min < 60) return `${min}分前`;
  const h = Math.floor(min / 60);
  if (h < 24) return `${h}時間前`;
  const d = Math.floor(h / 24);
  if (d < 30) return `${d}日前`;
  return new Date(iso).toLocaleDateString('ja-JP');
}

interface Props {
  item: FeedItem;
  isRead: boolean;
  isFav: boolean;
  onRead: (id: string) => void;
  onToggleFav: (id: string) => void;
}

export default function ItemCard({ item, isRead, isFav, onRead, onToggleFav }: Props) {
  return (
    <li
      className={`rounded-lg border px-4 py-3 transition-colors ${
        isRead
          ? 'border-slate-800 bg-slate-900/40 opacity-50'
          : 'border-slate-800 bg-slate-900 hover:border-slate-600'
      }`}
    >
      <div className="flex items-start gap-2">
        <button
          onClick={() => onToggleFav(item.id)}
          title={isFav ? 'お気に入り解除' : 'お気に入り'}
          className={`mt-0.5 shrink-0 text-lg leading-none ${
            isFav ? 'text-amber-400' : 'text-slate-600 hover:text-amber-400'
          }`}
        >
          {isFav ? '★' : '☆'}
        </button>
        <div className="min-w-0 flex-1">
          <a
            href={item.url}
            target="_blank"
            rel="noopener noreferrer"
            onClick={() => onRead(item.id)}
            onAuxClick={() => onRead(item.id)}
            className={`block font-medium leading-snug hover:underline ${
              isRead ? 'text-slate-400' : 'text-slate-100'
            }`}
          >
            {item.title}
          </a>
          {item.summary && (
            <p className="mt-1 line-clamp-2 text-xs text-slate-500">{item.summary}</p>
          )}
          <div className="mt-1.5 flex flex-wrap items-center gap-2 text-xs text-slate-500">
            <span className="rounded bg-slate-800 px-1.5 py-0.5 text-slate-400">
              {item.category}
            </span>
            <span>{item.source}</span>
            <span>{timeAgo(item.published)}</span>
          </div>
        </div>
      </div>
    </li>
  );
}
