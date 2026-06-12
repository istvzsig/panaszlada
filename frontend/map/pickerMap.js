import { state } from "../state/state.js";

export function initPickerMap() {
  if (state.pickerMap) return;

  state.pickerMap = new maplibregl.Map({
    container: "picker-map",
    style:
      "https://tiles.basemaps.cartocdn.com/gl/dark-matter-gl-style/style.json",
    center: [19.04, 47.49],
    zoom: 6,
  });

  state.pickerMap.on("click", (e) => {
    state.selectedLng = e.lngLat.lng;
    state.selectedLat = e.lngLat.lat;

    if (state.pickerMarker) state.pickerMarker.remove();

    state.pickerMarker = new maplibregl.Marker({ color: "#ef4444" })
      .setLngLat([state.selectedLng, state.selectedLat])
      .addTo(state.pickerMap);
  });
}
