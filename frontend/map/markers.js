import { popupHTML } from "../ui/popup.js";
import { categoryIcon } from "../utils/constants.js";
import { state } from "../state/state.js";

export function clearMarkers() {
  state.markers.forEach((m) => m.remove());
  state.markers = [];
}

export function renderMarkers(reports, map) {
  clearMarkers();

  reports.forEach((r) => {
    const icon = categoryIcon[r.category || "other"] || "📍";

    const el = document.createElement("div");
    el.style.fontSize = "22px";
    el.style.cursor = "pointer";
    el.textContent = icon;

    el.addEventListener("click", (e) => {
      e.stopPropagation();

      new maplibregl.Popup()
        .setLngLat([r.longitude, r.latitude])
        .setHTML(popupHTML(r))
        .addTo(map);
    });

    const marker = new maplibregl.Marker({
      element: el,
    })
      .setLngLat([r.longitude, r.latitude])
      .addTo(map);

    state.markers.push(marker);
  });
}
