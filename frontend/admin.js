const API =
  window.location.hostname === "localhost"
    ? "http://localhost:8080/api/v1"
    : "https://panaszlada.onrender.com/api/v1";

const API_KEY = "goibbTBMbATs65H3sdkgP4TYgw6DzFLZsATzOR5FzHI";

async function load() {
  const res = await fetch(API + "/reports", {
    headers: {
      "X-API-KEY": API_KEY,
    },
  });

  const data = await res.json();
  render(data);
}

function render(data) {
  if (data.length > 0) {
    document.getElementById("list").innerHTML = data
      .map(
        (r) => `
      <div class="card">

        <b>${r.title}</b>

        <div class="small">${r.description}</div>

        <div style="margin-top:8px">
          <span class="badge ${r.status}">
            ${r.status}
          </span>
        </div>

        <div class="small">
          ${r.tracking_code}
        </div>

        <div class="actions">

          <select onchange="updateStatus('${r.id}', this.value)">
            ${["open", "in_progress", "resolved", "invalid"]
              .map(
                (s) => `
                <option value="${s}" ${s === r.status ? "selected" : ""}>
                  ${s}
                </option>
              `,
              )
              .join("")}
          </select>

          <button class="danger" onclick="deleteReport('${r.id}')">
            Delete
          </button>

        </div>

      </div>
    `,
      )
      .join("");
  }
}

async function updateStatus(id, status) {
  await fetch(`${API}/reports/${id}`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
      "X-API-KEY": API_KEY,
    },
    body: JSON.stringify({ status }),
  });

  load();
}

async function deleteReport(id) {
  await fetch(`${API}/reports/${id}`, {
    method: "DELETE",
    headers: {
      "X-API-KEY": API_KEY,
    },
  });

  load();
}

load();
