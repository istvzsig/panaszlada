import {
  categoryIcon,
  categoryLabel,
  statusLabel,
} from "../utils/constants.js";

export function popupHTML(r) {
  const cat = r.category || "other";

  return `
    <strong>${r.title}</strong><br><br>
    ${r.description}<br><br>

    <div>
      ${categoryIcon[cat]} ${categoryLabel[cat]}
    </div>

    <div>
      <b>${statusLabel[r.status] || r.status}</b>
    </div>

    <div style="margin-top:6px">
      <small>${r.tracking_code}</small>
    </div>
  `;
}
