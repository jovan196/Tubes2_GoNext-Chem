import { useState } from "react";
import axios from "axios";
import RecipeTree from "./components/RecipeTree";

function App() {
  const [target, setTarget] = useState("");
  const [algorithm, setAlgorithm] = useState("BFS");
  const [multiple, setMultiple] = useState(false);
  const [maxRecipe, setMaxRecipe] = useState(1);
  const [result, setResult] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const payload = {
      target,
      algorithm,
      mode: multiple ? "multiple" : "single",
      max: multiple ? parseInt(maxRecipe) : 1,
    };

    try {
      const res = await axios.post("http://localhost:8080/api/search", payload);
      setResult(res.data);
    } catch (err) {
      alert("Gagal mengambil data dari backend.");
      console.error(err);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <form
        onSubmit={handleSubmit}
        className="max-w-xl mx-auto bg-white p-6 rounded-lg shadow space-y-4"
      >
        <h1 className="text-2xl font-bold text-center">Recipe Finder</h1>

        <input
          type="text"
          placeholder="Nama elemen (mis: Brick)"
          value={target}
          onChange={(e) => setTarget(e.target.value)}
          className="w-full p-2 border rounded"
          required
        />

        <div className="flex justify-around">
          <label className="flex items-center gap-2">
            <input
              type="radio"
              name="algorithm"
              value="BFS"
              checked={algorithm === "BFS"}
              onChange={(e) => setAlgorithm(e.target.value)}
            />
            BFS
          </label>
          <label className="flex items-center gap-2">
            <input
              type="radio"
              name="algorithm"
              value="DFS"
              checked={algorithm === "DFS"}
              onChange={(e) => setAlgorithm(e.target.value)}
            />
            DFS
          </label>
        </div>

        <label className="flex items-center gap-2">
          <input
            type="checkbox"
            checked={multiple}
            onChange={(e) => setMultiple(e.target.checked)}
          />
          Multiple Recipe Mode
        </label>

        {multiple && (
          <input
            type="number"
            min={1}
            placeholder="Max jumlah recipe"
            value={maxRecipe}
            onChange={(e) => setMaxRecipe(e.target.value)}
            className="w-full p-2 border rounded"
          />
        )}

        <button
          type="submit"
          className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
        >
          Cari Recipe
        </button>
      </form>

      {result && (
        <div className="mt-8">
          <h2 className="text-xl font-semibold text-center mb-4">
            Hasil Recipe
          </h2>
          <RecipeTree data={result} />
        </div>
      )}
    </div>
  );
}

export default App;
