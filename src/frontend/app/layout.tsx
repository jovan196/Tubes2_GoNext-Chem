// app/layout.tsx
import '../styles/globals.css';
import { ReactNode } from 'react';

export const metadata = {
  title: 'Little Alchemy 2 Calculator',
  description: 'BFS/DFS recipe visualizer for Little Alchemy 2',
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <body className="bg-gray-100">
        <header className="bg-blue-600 text-white py-4 shadow">
          <h1 className="text-center text-2xl font-bold">
            Little Alchemy 2 Calculator
          </h1>
        </header>
        <main className="container mx-auto p-4">{children}</main>
      </body>
    </html>
  );
}
