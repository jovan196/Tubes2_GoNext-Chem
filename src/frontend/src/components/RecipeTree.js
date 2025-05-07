function RecipeTree({ data }) {
    if (!data || data.length === 0) return <p className="text-center">Tidak ada recipe ditemukan.</p>;
  
    return (
      <div className="space-y-6">
        {data.map((tree, idx) => (
          <div key={idx} className="border p-4 rounded bg-white shadow">
            <p className="font-bold mb-2">Recipe {idx + 1}:</p>
            <pre className="whitespace-pre-wrap">{JSON.stringify(tree, null, 2)}</pre>
          </div>
        ))}
      </div>
    );
  }
  
  export default RecipeTree;
  