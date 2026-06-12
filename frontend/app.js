import { initPublicMap, loadReports } from "./map/publicMap.js";
import { initPickerMap } from "./map/pickerMap.js";
import { createReport } from "./service/reportService.js";
import { state } from "./state/state.js";
import { setStatus } from "./ui/toast.js";
import { TITLE_LIMIT, DESC_LIMIT, HUNGARY_BOUNDS } from "./utils/constants.js";
import { isInHungary } from "./utils/validation.js";

import { setStatusFilter, applyFilters } from "./ui/filters.js";
import { renderList } from "./ui/list.js";
import { renderStatusSummary } from "./ui/status.js";
import { closeModal } from "./ui/modal.js";

window.setStatusFilter = setStatusFilter;

window.focusReport = (code) => {
  const r = state.reports.find((x) => x.tracking_code === code);
  if (!r) return;

  state.publicMap.flyTo({
    center: [r.longitude, r.latitude],
    zoom: 15,
  });
};

function initForm() {
  const form = document.getElementById("report-form");

  form.addEventListener("submit", async (e) => {
    e.preventDefault();

    const title = document.getElementById("title").value;
    const desc = document.getElementById("desc").value;
    const category = document.getElementById("category")?.value || "other";

    if (!title || title.length > TITLE_LIMIT)
      return setStatus("Hibás cím", "error");

    if (!desc || desc.length > DESC_LIMIT)
      return setStatus("Hibás leírás", "error");

    if (state.selectedLat === null || state.selectedLng === null)
      return setStatus("Válassz helyet!", "error");

    if (!isInHungary(state.selectedLng, state.selectedLat, HUNGARY_BOUNDS))
      return setStatus("Csak Magyarországon!", "error");

    await createReport({
      title,
      description: desc,
      category,
      latitude: state.selectedLat,
      longitude: state.selectedLng,
    });

    setStatus("Sikeres bejelentés", "success");
  });
}

function main() {
  initPublicMap();

  document.getElementById("new-report-btn").onclick = () => {
    document.getElementById("modal").classList.add("show");
    initPickerMap();
  };

  document.getElementById("close-modal-btn").onclick = () => {
    closeModal();
  };

  initForm();
  loadReports();
}

main();
