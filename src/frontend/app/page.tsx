// app/page.tsx
'use client';

import { useState, FormEvent } from 'react';
import SearchMenu from '../components/SearchMenu';
import RecipeTree, { SearchResponse } from '../components/RecipeTree';

export default function HomePage() {
  const [target, setTarget] = useState('');
  const [algorithm, setAlgorithm] = useState<'bfs' | 'dfs'>('bfs');
  const [mode, setMode] = useState<'single' | 'multiple'>('single');
  const [maxResults, setMaxResults] = useState(3);
  const [result, setResult] = useState<SearchResponse[] | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!target.trim()) {
      setError('Masukkan elemen target.');
      return;
    }
    setLoading(true);
    setError(null);
    setResult(null);
    try {
      const res = await fetch('/api/search', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ target, algorithm, mode, max: maxResults }),
      });
      if (!res.ok) throw new Error(await res.text());
      const data: SearchResponse[] = await res.json();
      setResult(data);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      <SearchMenu
        target={target}
        setTarget={setTarget}
        algorithm={algorithm}
        setAlgorithm={setAlgorithm}
        mode={mode}
        setMode={setMode}
        maxResults={maxResults}
        setMaxResults={setMaxResults}
        loading={loading}
        error={error}
        onSubmit={handleSubmit}
      />

      {result && result.length > 0 && (
        <div className="space-y-8">
          {result.map((resp, i) => (
            <div key={i}>
              <h2 className="font-semibold mb-2">
                Solution {i + 1} for "{resp.result}"
              </h2>
              <RecipeTree response={resp} />
            </div>
          ))}
        </div>
      )}

      {result && result.length === 0 && (
        <p className="text-center text-gray-500">
          Element not reachable.
        </p>
      )}
    </div>
  );
}
