import { state } from "../state/state.js";
import { calculateStats } from "./stats.js";

export function renderStatusSummary() {
  const container = document.getElementById("status");
  if (!container) return;

  const stats = calculateStats(state.reports);

  container.innerHTML = `
    <div class="status-summary">

      <div class="status-box all" onclick="window.setStatusFilter('all')">
        <strong>${state.reports.length}</strong>
        <span>Összes</span>
      </div>

      <div class="status-box open" onclick="window.setStatusFilter('open')">
        <strong>${stats.open}</strong>
        <span>Nyitott</span>
      </div>

      <div class="status-box in_progress" onclick="window.setStatusFilter('in_progress')">
        <strong>${stats.in_progress}</strong>
        <span>Folyamatban</span>
      </div>

      <div class="status-box resolved" onclick="window.setStatusFilter('resolved')">
        <strong>${stats.resolved}</strong>
        <span>Megoldott</span>
      </div>

      <div class="status-box invalid" onclick="window.setStatusFilter('invalid')">
        <strong>${stats.invalid}</strong>
        <span>Érvénytelen</span>
      </div>

    </div>
  `;
}
