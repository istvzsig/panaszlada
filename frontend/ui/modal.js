import { state } from "../state/state.js";

export function closeModal() {
  const modal = document.getElementById("modal");
  modal.classList.remove("show");

  // reset form state
  const form = document.getElementById("report-form");
  if (form) form.reset();

  // reset selected coords
  state.selectedLat = null;
  state.selectedLng = null;

  // remove picker marker if exists
  if (state.pickerMarker) {
    state.pickerMarker.remove();
    state.pickerMarker = null;
  }

  // optional: reset map view
  if (state.pickerMap) {
    state.pickerMap.flyTo({
      center: [19.04, 47.49],
      zoom: 6,
    });
  }
}
