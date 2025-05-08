import React, { useState } from "react";
import SearchForm from "./components/SearchForm";
import RecipeTree from "./components/RecipeTree"; // pastikan ini ditambahkan

function App() {
  const [result, setResult] = useState(null);

  return (
    <div className="min-h-screen p-8 bg-gray-50">
      <h1 className="text-2xl font-bold mb-4">Recipe Finder</h1>
      <SearchForm setResult={setResult} />

      {result && (
        <div className="mt-8">
          <h2 className="font-bold mb-4">Hasil:</h2>
          <RecipeTree data={result} />
        </div>
      )}
    </div>
  );
}

export default App;
