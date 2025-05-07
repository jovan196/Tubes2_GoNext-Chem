import React, { useState } from "react";
import SearchForm from "./components/SearchForm";

function App() {
  const [result, setResult] = useState(null);

  return (
    <div className="min-h-screen p-8 bg-gray-50">
      <h1 className="text-2xl font-bold mb-4">Recipe Finder</h1>
      <SearchForm setResult={setResult} />

      {result && (
        <div className="mt-8">
          <h2 className="font-bold">Hasil:</h2>
          <pre className="bg-white p-4 rounded shadow-md">
            {JSON.stringify(result, null, 2)}
          </pre>
        </div>
      )}
    </div>
  );
}

export default App;
