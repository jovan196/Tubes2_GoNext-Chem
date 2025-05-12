import React, { useState } from "react";
import axios from "axios";

function SearchForm({ setResult }) {
  const [target, setTarget] = useState("");
  const [algorithm, setAlgorithm] = useState("BFS");
  const [multiple, setMultiple] = useState(false);
  const [maxRecipe, setMaxRecipe] = useState(1);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    const payload = {
      target: target.trim(),
      algorithm,
      mode: multiple ? "multiple" : "single",
      max: multiple ? parseInt(maxRecipe, 10) : 1,
    };

    axios.post(
        "http://localhost:8080/api/search",
        payload
      ).then(
        res => {
          setResult(res.data);
        }
      ).catch(error => {
        let errorMessage = error.message
        if (error.response && error.response.data) {
          errorMessage = error.response.data.message || errorMessage
        }
        alert(errorMessage);
      }).finally(() => setLoading(false));

    // try {
    //   const res = await axios.post(
    //     "http://localhost:8080/api/search",
    //     payload
    //   );
    //   // res.data is array of { result, steps, timeMs, visitedCount }
    //   setResult(res.data);
    // } catch (err) {
    //   //alert("Gagal mengambil data dari backend.");
    //   alert(err);
    //   console.error(err);
    // } finally {
    //   setLoading(false);
    // }
  };

  return (
    <form
      onSubmit={handleSubmit}
      className="max-w-xl mx-auto bg-white p-6 rounded-lg shadow space-y-4"
    >
      <input
        type="text"
        placeholder="Nama elemen (mis: Brick)"
        value={target}
        onChange={(e) => setTarget(e.target.value)}
        className="w-full p-2 border rounded"
        required
      />

      <div className="flex justify-around">
        {['BFS', 'DFS'].map((algo) => (
          <label key={algo} className="flex items-center gap-2">
            <input
              type="radio"
              name="algorithm"
              value={algo}
              checked={algorithm === algo}
              onChange={(e) => setAlgorithm(e.target.value)}
            />
            {algo}
          </label>
        ))}
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
          required
        />
      )}

      <button
        type="submit"
        disabled={loading}
        className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 w-full disabled:opacity-50"
      >
        {loading ? 'Mencari...' : 'Cari Recipe'}
      </button>
    </form>
  );
}

export default SearchForm;