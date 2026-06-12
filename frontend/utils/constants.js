export const API =
  window.location.hostname === "localhost"
    ? "https://panaszlada.onrender.com/api/v1"
    : "https://panaszlada.onrender.com/api/v1";

export const HUNGARY_BOUNDS = [
  [16.0, 45.7],
  [23.1, 48.7],
];

export const TITLE_LIMIT = 80;
export const DESC_LIMIT = 1000;

export const categoryIcon = {
  safety: "⚠️",
  infrastructure: "🏗️",
  homeless: "🧍",
  other: "📍",
};

export const categoryLabel = {
  safety: "Biztonság",
  infrastructure: "Infrastruktúra",
  homeless: "Hajléktalan",
  other: "Egyéb",
};

export const statusLabel = {
  open: "Nyitott",
  in_progress: "Folyamatban",
  resolved: "Megoldott",
  invalid: "Érvénytelen",
};
