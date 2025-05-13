import React, { useMemo } from "react";
import Tree from "react-d3-tree";

const BASIC = new Set(["Air", "Water", "Earth", "Fire", "Time"]);

const RecipeTree = ({ data }) => {
  const treeDatas = useMemo(() => {
    if (!Array.isArray(data) || data.length === 0) return [];
    return data.map((resp) => resp.tree); // pakai .tree langsung dari backend
  }, [data]);

  if (!treeDatas.length) {
    return <p className="text-center">Tidak ada recipe ditemukan.</p>;
  }

  return (
    <div className="space-y-8">
      {data.map((resp, idx) => {
        const { result, timeMs, visitedCount } = resp;
        const treeData = treeDatas[idx];

        return (
          <div key={idx} className="border rounded p-4 bg-gray-50">
            <h2 className="font-bold mb-2">
              Recipe {idx + 1} untuk <em>{result}</em>
            </h2>

            {treeData ? (
              <>
                <p className="text-sm text-gray-600 mb-4">
                  Waktu: {timeMs} ms &middot; Simpul Dikunjungi: {visitedCount}
                </p>
                <div className="w-full h-[400px]">
                  <Tree
                    data={treeData}
                    orientation="vertical"
                    pathFunc="elbow"
                    translate={{ x: 300, y: 100 }}
                    collapsible={false}
                    renderCustomNodeElement={({ nodeDatum }) => {
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
                            strokeWidth={1}
                            rx={6}
                          />
                          <text
                            fill="#111"
                            x={0}
                            y={5}
                            textAnchor="middle"
                            fontSize={12}
                            fontWeight="normal"
                          >
                            {nodeDatum.name}
                          </text>
                        </g>
                      );
                    }}
                  />
                </div>
              </>
            ) : (
              <p className="text-sm text-red-500">
                Tidak ada recipe yang bisa dibentuk.
              </p>
            )}
          </div>
        );
      })}
    </div>
  );
};

export default RecipeTree;
