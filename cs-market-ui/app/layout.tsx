import './globals.css';
import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import { Providers } from './providers';
import { Navbar } from '@/components/navbar';

// const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'CS2 Skins Marketplace',
  description: 'Buy and sell CS2 skins with real-time prices and secure trading',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={`min-h-screen bg-background`}>
        <Providers>
          <Navbar />
          <main className="container mx-auto px-4 py-6">
            {children}
          </main>
        </Providers>
      </body>
    </html>
  );
}