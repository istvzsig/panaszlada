import { state } from "../state/state.js";
import { categoryIcon, categoryLabel } from "../utils/constants.js";
import { escapeHtml } from "../utils/escape.js";

export function renderList(containerId = "list") {
  const container = document.getElementById(containerId);
  if (!container) return;

  const data = state.filteredReports || state.reports || [];

  container.innerHTML = data
    .map((r) => {
      const cat = r.category || "other";

      const icon = categoryIcon[cat] || "📍";
      const label = categoryLabel[cat] || cat;

      return `
        <div class="card" data-code="${r.tracking_code}">
          <b>${escapeHtml(r.title)}</b>
          <div class="small">${escapeHtml(r.description)}</div>

          <div class="badge ${r.status}">
            ${r.status}
          </div>

          <div class="small">
            ${icon} ${label}
          </div>

          <div class="small">${r.tracking_code}</div>
        </div>
      `;
    })
    .join("");

  // safer event binding instead of onclick
  container.querySelectorAll(".card").forEach((el) => {
    el.addEventListener("click", () => {
      const code = el.getAttribute("data-code");
      window.focusReport?.(code);
    });
  });
}
