import { state } from "../state/state.js";
import { renderList } from "./list.js";

export function setStatusFilter(status) {
  if (status === "all") {
    state.activeStatusFilter = null;
  } else {
    state.activeStatusFilter =
      state.activeStatusFilter === status ? null : status;
  }

  applyFilters();
}

export function applyFilters() {
  const data = state.reports;

  if (state.activeStatusFilter) {
    state.filteredReports = data.filter(
      (r) => r.status === state.activeStatusFilter,
    );
  } else {
    state.filteredReports = data;
  }

  renderList();
}
