// components/RecipeTree.tsx
import React, { useMemo, useRef, useLayoutEffect, useState } from 'react';
import Tree from 'react-d3-tree';

const BASIC = new Set(['Air', 'Water', 'Earth', 'Fire', 'Time']);

export interface RecipeStep {
  product: string;
  ingredients: string[];
}
export interface SearchResponse {
  result: string;
  steps?: RecipeStep[];      // ←  tanda tanya = bisa tidak ada
}
interface Props {
  response: SearchResponse;
}

const RecipeTree: React.FC<Props> = ({ response }) => {
  const { result, steps } = response;

  const buildTree = (root: string, s?: RecipeStep[]) => {
    // jika tidak ada langkah → node tunggal
    if (!s || s.length === 0) return { name: root };

    const map = new Map<string, string[]>();
    s.forEach((st) => map.set(st.product, st.ingredients));

    const visited = new Set<string>();
    const recurse = (el: string): any => {
      if (visited.has(el) || BASIC.has(el) || !map.has(el)) {
        return { name: el };
      }
      visited.add(el);
      return {
        name: el,
        children: (map.get(el) || []).map(recurse),
      };
    };
    return recurse(root);
  };

  const treeData = useMemo<any>(() => buildTree(result, steps), [
    result,
    steps,
  ]);

  // ---------------------- posisi & render node ----------------------
  const containerRef = useRef<HTMLDivElement>(null);
  const [translate, setTranslate] = useState({ x: 0, y: 50 });

  useLayoutEffect(() => {
    if (containerRef.current) {
      const w = containerRef.current.getBoundingClientRect().width;
      setTranslate({ x: w / 2, y: 50 });
    }
  }, []);

  const renderNode = ({ nodeDatum }: any) => {
    const name = nodeDatum.name as string;
    const fill = BASIC.has(name) ? '#D1FAE5' : '#DBEAFE';
    return (
      <g>
        <rect x={-50} y={-20} width={100} height={40} rx={6}
              fill={fill} stroke="#555" strokeWidth={1.5} />
        <text x={0} y={5} textAnchor="middle" fontSize={12} fontWeight="bold">
          {name}
        </text>
      </g>
    );
  };

  return (
    <div
      ref={containerRef}
      className="w-full h-[600px] overflow-auto border rounded bg-white shadow"
    >
      <Tree
        data={[treeData] as any}
        orientation="vertical"
        translate={translate}
        pathFunc="elbow"
        collapsible={false}
        renderCustomNodeElement={renderNode}
      />
    </div>
  );
};

export default RecipeTree;
