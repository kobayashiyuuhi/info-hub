export interface FeedItem {
  id: string;
  title: string;
  url: string;
  source: string;
  category: string;
  published: string; // RFC3339
  summary?: string;
}

export interface FeedStatus {
  name: string;
  url: string;
  category: string;
  status: 'ok' | 'error';
  error?: string;
  itemCount: number;
}

export interface Meta {
  generatedAt: string;
  feeds: FeedStatus[];
}
