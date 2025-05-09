import { FormEvent } from 'react';

export type SearchMenuProps = {
  target: string;
  setTarget: (v: string) => void;
  algorithm: 'bfs' | 'dfs';
  setAlgorithm: (a: 'bfs' | 'dfs') => void;
  mode: 'single' | 'multiple';
  setMode: (m: 'single' | 'multiple') => void;
  maxResults: number;
  setMaxResults: (n: number) => void;
  loading: boolean;
  error: string | null;
  onSubmit: (e: FormEvent) => void;
};

export default function SearchMenu({
  target, setTarget,
  algorithm, setAlgorithm,
  mode, setMode,
  maxResults, setMaxResults,
  loading, error,
  onSubmit
}: SearchMenuProps) {
  return (
    <>
      <form onSubmit={onSubmit} className="bg-white p-6 rounded shadow space-y-4">
        <div>
          <label className="block font-medium mb-1">Target Element</label>
          <input
            type="text"
            value={target}
            onChange={e => setTarget(e.target.value)}
            className="w-full border rounded p-2"
            placeholder="e.g. Fire"
          />
        </div>
        <div className="flex items-center space-x-6">
          <div>
            <span className="font-medium">Algorithm:</span>
            <label className="ml-2">
              <input
                type="radio"
                name="algo"
                checked={algorithm === 'bfs'}
                onChange={() => setAlgorithm('bfs')}
                className="mr-1"
              />
              BFS
            </label>
            <label className="ml-2">
              <input
                type="radio"
                name="algo"
                checked={algorithm === 'dfs'}
                onChange={() => setAlgorithm('dfs')}
                className="mr-1"
              />
              DFS
            </label>
          </div>
          <div>
            <span className="font-medium">Mode:</span>
            <label className="ml-2">
              <input
                type="radio"
                name="mode"
                checked={mode === 'single'}
                onChange={() => setMode('single')}
                className="mr-1"
              />
              Single
            </label>
            <label className="ml-2">
              <input
                type="radio"
                name="mode"
                checked={mode === 'multiple'}
                onChange={() => setMode('multiple')}
                className="mr-1"
              />
              Multiple
            </label>
            {mode === 'multiple' && (
              <input
                type="number"
                min={1}
                value={maxResults}
                onChange={e => setMaxResults(+e.target.value)}
                className="ml-2 w-16 border rounded p-1"
              />
            )}
          </div>
        </div>
        <button
          type="submit"
          disabled={loading}
          className="w-full bg-blue-600 text-white py-2 rounded"
        >
          {loading ? 'Loading…' : 'Search'}
        </button>
      </form>
      {loading && <p className="text-center">Memproses…</p>}
      {error && <p className="text-red-600 text-center">{error}</p>}
    </>
  );
}
