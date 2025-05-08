import React, { useMemo } from "react";
import Tree from "react-d3-tree";

// Elemen dasar termasuk Time
const BASIC = new Set(["Air", "Water", "Earth", "Fire", "Time"]);

const RecipeTree = ({ data }) => {
  const steps = data?.[0]?.steps || [];
  const result = data?.[0]?.result || "";

  // Build peta dari product -> ingredients
  const buildTreeData = (result, steps) => {
    const recipeMap = new Map();
    for (const step of steps) {
      recipeMap.set(step.product, step.ingredients);
    }
  
    const visited = new Set(); // untuk deteksi siklus
  
    const build = (element) => {
      if (visited.has(element)) {
        return { name: element }; // cegah infinite loop
      }
  
      visited.add(element);
  
      const ingredients = recipeMap.get(element);
      if (!ingredients) {
        return { name: element };
      }
  
      const children = ingredients.map(build);
      return {
        name: element,
        children,
      };
    };
  
    return build(result);
  };
  
  
  

  const treeData = useMemo(() => {
    if (!steps.length || !result) return null;
    return buildTreeData(result, steps);
  }, [result, steps]);

  const renderCustomNode = ({ nodeDatum }) => {
    const isBasic = BASIC.has(nodeDatum.name);
    const color = isBasic ? "#D1FAE5" : "#DBEAFE";

    return (
      <g>
        <rect
          width={80}
          height={30}
          x={-40}
          y={-15}
          fill={color}
          stroke="#555"
          strokeWidth={1.5}
          rx={6}
        />
        <text
          fill="#111"
          stroke="none"
          x={0}
          y={5}
          textAnchor="middle"
          fontSize={12}
          fontWeight="bold"
        >
          {nodeDatum.name}
        </text>
      </g>
    );
  };

  if (!treeData) {
    return <p className="text-center">Tidak ada recipe ditemukan.</p>;
  }

  return (
    <div className="w-full h-[600px] border rounded bg-white shadow overflow-auto">
      <Tree
        data={treeData}
        orientation="vertical"
        translate={{ x: 400, y: 100 }}
        pathFunc="elbow"
        collapsible={false}
        renderCustomNodeElement={renderCustomNode}
      />
    </div>
  );
};

export default RecipeTree;
