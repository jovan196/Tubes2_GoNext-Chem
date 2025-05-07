<<<<<<< HEAD
import pandas as pd
from collections import defaultdict
import json
import csv

def scrape_with_pandas():
    url = "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"
    tables = pd.read_html(url)

    graph = defaultdict(list)
    basic = {"Air","Earth","Fire","Water"}

    for df in tables:
        cols = [c.lower() for c in df.columns]
        if "element" not in cols or not any(c in cols for c in ("combinations","recipes")):
            continue

        name_col = cols.index("element")
        combo_col = cols.index("combinations") if "combinations" in cols else cols.index("recipes")

        for _, row in df.iterrows():
            name = row[name_col]
            raw = str(row[combo_col])
            # raw might look like "Dough + Fruit, Dough + Dough" etc.
            for part in raw.replace("\n",", ").split(","):
                if "+" not in part:
                    continue
                a,b = [p.strip() for p in part.split("+",1)]
                graph[name].append((a,b))

    return graph

if __name__ == "__main__":
    g = scrape_with_pandas()
    print(f"Found {len(g)} elements with recipes.")

    # save JSON
    with open("elements_graph.json","w",encoding="utf-8") as jf:
        json.dump(g, jf, default=list, indent=2)

    # save CSV edge‑list
    with open("elements_graph.csv","w",newline="",encoding="utf-8") as cf:
        writer = csv.writer(cf)
        writer.writerow(["Product","Ingredient1","Ingredient2"])
        for prod, recs in g.items():
            for a,b in recs:
                writer.writerow([prod,a,b])
=======
import pandas as pd
from collections import defaultdict
import json
import csv

def scrape_with_pandas():
    url = "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"
    tables = pd.read_html(url)

    graph = defaultdict(list)
    basic = {"Air","Earth","Fire","Water"}

    for df in tables:
        cols = [c.lower() for c in df.columns]
        if "element" not in cols or not any(c in cols for c in ("combinations","recipes")):
            continue

        name_col = cols.index("element")
        combo_col = cols.index("combinations") if "combinations" in cols else cols.index("recipes")

        for _, row in df.iterrows():
            name = row[name_col]
            raw = str(row[combo_col])
            # raw might look like "Dough + Fruit, Dough + Dough" etc.
            for part in raw.replace("\n",", ").split(","):
                if "+" not in part:
                    continue
                a,b = [p.strip() for p in part.split("+",1)]
                graph[name].append((a,b))

    return graph

if __name__ == "__main__":
    g = scrape_with_pandas()
    print(f"Found {len(g)} elements with recipes.")

    # save JSON
    with open("elements_graph.json","w",encoding="utf-8") as jf:
        json.dump(g, jf, default=list, indent=2)

    # save CSV edge‑list
    with open("elements_graph.csv","w",newline="",encoding="utf-8") as cf:
        writer = csv.writer(cf)
        writer.writerow(["Product","Ingredient1","Ingredient2"])
        for prod, recs in g.items():
            for a,b in recs:
                writer.writerow([prod,a,b])
>>>>>>> c667bd805e708bd1eff2b7921b6fe17486ec371a
