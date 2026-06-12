import { renderMarkers } from "./markers.js";
import { renderList } from "../ui/list.js";
import { renderStatusSummary } from "../ui/status.js";
import { state } from "../state/state.js";
import { getReports } from "../service/reportService.js";
import { setStatus } from "../ui/toast.js";
import { API, HUNGARY_BOUNDS } from "../utils/constants.js";
import { applyFilters } from "../ui/filters.js";

export function initPublicMap() {
  state.publicMap = new maplibregl.Map({
    container: "public-map",
    style:
      "https://tiles.basemaps.cartocdn.com/gl/dark-matter-gl-style/style.json",
    center: [19.04, 47.49],
    zoom: 6,
    maxBounds: HUNGARY_BOUNDS,
  });

  state.publicMap.addControl(new maplibregl.NavigationControl());

  state.publicMap.on("load", async () => {
    await loadReports();
  });
}

export async function loadReports() {
  try {
    const res = await fetch(API + "/reports");
    const data = await res.json();

    state.reports = Array.isArray(data) ? data : [];

    state.filteredReports = state.reports;

    renderMarkers(state.reports, state.publicMap);
    renderList();
    renderStatusSummary();
    applyFilters();
  } catch (e) {
    console.error(e);
    setStatus("Hiba a betöltésnél", "error");
  }
}
