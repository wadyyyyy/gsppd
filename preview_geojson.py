import json
from pathlib import Path

import matplotlib.pyplot as plt
from matplotlib.patches import PathPatch
from matplotlib.path import Path as MPath

filepath = Path("data/areas.geojson")

with open(filepath, "r") as f:
    data = json.load(f)

fig, ax = plt.subplots(figsize=(8, 8))

colors = plt.rcParams["axes.prop_cycle"].by_key()["color"]

for index, feature in enumerate(data["features"]):
    geom = feature["geometry"]
    coords = geom["coordinates"]

    color = colors[index % len(colors)]

    polys = coords if geom["type"] == "MultiPolygon" else [coords]

    for poly in polys:
        all_path_codes = []
        all_path_verts = [] 

        for ring in poly:
            verts = [(p[0], p[1]) for p in ring]
            all_path_verts.extend(verts)
            codes = [MPath.LINETO] * len(verts)
            codes[0] = MPath.MOVETO
            all_path_codes.extend(codes)

        path = MPath(all_path_verts, all_path_codes)
        patch = PathPatch(path, facecolor=color, edgecolor="black", lw=0.5, alpha=0.6)
        ax.add_patch(patch)

ax.autoscale_view()
ax.set_aspect("equal")
plt.title("Визуализация с учетом вырезанных зон")
plt.show()
