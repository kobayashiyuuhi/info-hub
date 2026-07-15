import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'info-hub',
  description: 'Personal RSS / X information dashboard',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ja">
      <body className="min-h-screen antialiased">{children}</body>
    </html>
  );
}
