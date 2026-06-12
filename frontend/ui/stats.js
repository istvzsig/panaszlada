import { state } from "../state/state.js";

export function calculateStats(data = []) {
  const stats = {
    open: 0,
    in_progress: 0,
    resolved: 0,
    invalid: 0,
  };

  data.forEach((r) => {
    if (r?.status && stats[r.status] !== undefined) {
      stats[r.status]++;
    }
  });

  return stats;
}
