import { API } from "../utils/constants.js";

export async function getReports() {
  const res = await fetch(API + "/reports");
  return res.json();
}

export async function createReport(payload) {
  const res = await fetch(API + "/reports", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  if (!res.ok) throw new Error("failed");
  return res.json();
}
