const API =
  window.location.hostname === "localhost"
    ? "http://localhost:8080/api/v1"
    : "https://panaszlada.onrender.com/api/v1";

// =====================
// CONFIG
// =====================
const HUNGARY_BOUNDS = [
  [16.0, 45.7],
  [23.1, 48.7],
];

const categoryIcon = {
  safety: "⚠️",
  infrastructure: "🏗️",
  homeless: "🧍",
  other: "📍",
};

const categoryLabel = {
  safety: "Biztonság",
  infrastructure: "Infrastruktúra",
  homeless: "Hajléktalan",
  other: "Egyéb",
};

const TITLE_LIMIT = 80;
const DESC_LIMIT = 1000;

const reportTitle = document.getElementById("title");
const reportDescription = document.getElementById("desc");
const submitBtn = document.getElementById("submit-btn");

// =====================
// STATS
// =====================
const stats = {
  open: 0,
  in_progress: 0,
  resolved: 0,
  invalid: 0,
};

// =====================
// STATE
// =====================
let selectedLat = null;
let selectedLng = null;

let publicMap = null;
let pickerMap = null;
let pickerMarker = null;

let reportData = [];

// marker storage
window.__markers = [];

// =====================
// STATUS
// =====================
let statusTimeout = null;
let activeStatusFilter = null;

function calculateStats(data) {
  const stats = {
    open: 0,
    in_progress: 0,
    resolved: 0,
    invalid: 0,
  };

  data.forEach((r) => {
    if (stats[r.status] !== undefined) {
      stats[r.status]++;
    }
  });

  return stats;
}

function setStatus(message, type = "info") {
  const container = document.getElementById("toast-container");

  const icons = {
    success: "✓",
    error: "✕",
    info: "ℹ",
  };

  const toast = document.createElement("div");

  toast.className = `toast ${type}`;

  toast.innerHTML = `
    <div class="toast-icon">
      ${icons[type] || "ℹ"}
    </div>

    <div class="toast-message">
      ${message}
    </div>
  `;

  container.appendChild(toast);

  setTimeout(() => {
    toast.remove();
  }, 4000);
}

// =====================
// VALIDATION
// =====================
function isInHungary(lng, lat) {
  return (
    lng >= HUNGARY_BOUNDS[0][0] &&
    lng <= HUNGARY_BOUNDS[1][0] &&
    lat >= HUNGARY_BOUNDS[0][1] &&
    lat <= HUNGARY_BOUNDS[1][1]
  );
}

// =====================
// API
// =====================
async function loadReports() {
  try {
    const res = await fetch(API + "/reports");
    const data = await res.json();

    reportData = data || [];

    // renderList(reportData);
    // renderMarkers(reportData);
    renderStatusSummary(reportData);
    // updateClusterData(reportData);
    applyFilters();
  } catch (e) {
    setStatus("Hiba a betöltésnél", "error");
  }
}

// =====================
// CREATE REPORT
// =====================
async function createReport() {
  const title = reportTitle.value;
  const desc = document.getElementById("desc").value;
  const category = document.getElementById("category")?.value || "other";

  if (!title || title.length > TITLE_LIMIT) {
    setStatus("Hibás cím", "error");
    return;
  }

  if (!desc || desc.length > DESC_LIMIT) {
    setStatus("Hibás leírás", "error");
    return;
  }

  if (selectedLat === null || selectedLng === null) {
    setStatus("Válassz helyet!", "error");
    return;
  }

  if (!isInHungary(selectedLng, selectedLat)) {
    setStatus("Csak Magyarországon lehet!", "error");
    return;
  }

  submitBtn.disabled = true;
  submitBtn.innerText = "Mentés...";

  try {
    const res = await fetch(API + "/reports", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        title,
        description: desc,
        category,
        latitude: selectedLat,
        longitude: selectedLng,
      }),
    });

    if (!res.ok) throw new Error("failed");

    setStatus(`Sikeres bejelentés: ${title}`, "success");

    closeModal();
    await loadReports();
  } catch (e) {
    setStatus("Hiba a küldésnél", "error");
  }

  submitBtn.disabled = false;
  submitBtn.innerText = "Mentés";
}

// =====================
// MAP INIT
// =====================
function initPublicMap() {
  publicMap = new maplibregl.Map({
    container: "public-map",
    style:
      "https://tiles.basemaps.cartocdn.com/gl/dark-matter-gl-style/style.json",
    center: [19.04, 47.49],
    zoom: 6,
    maxBounds: HUNGARY_BOUNDS,
  });

  publicMap.addControl(new maplibregl.NavigationControl());

  publicMap.on("load", () => {
    setupCluster();
    loadReports();
  });
}

// =====================
// CLUSTER (ONLY CLUSTER HERE)
// =====================
function setupCluster() {
  publicMap.addSource("reports", {
    type: "geojson",
    data: { type: "FeatureCollection", features: [] },
    cluster: true,
    clusterMaxZoom: 11,
    clusterRadius: 60,
  });

  publicMap.addLayer({
    id: "clusters",
    type: "circle",
    source: "reports",
    filter: ["has", "point_count"],
    paint: {
      "circle-color": "#3b82f6",
      "circle-radius": 18,
      "circle-opacity": 0.6,
    },
  });

  publicMap.addLayer({
    id: "cluster-count",
    type: "symbol",
    source: "reports",
    filter: ["has", "point_count"],
    layout: {
      "text-field": "{point_count_abbreviated}",
      "text-size": 12,
    },
  });

  publicMap.on("click", "clusters", (e) => {
    const features = publicMap.queryRenderedFeatures(e.point, {
      layers: ["clusters"],
    });

    const clusterId = features[0].properties.cluster_id;

    publicMap
      .getSource("reports")
      .getClusterExpansionZoom(clusterId, (err, zoom) => {
        if (err) return;

        publicMap.easeTo({
          center: features[0].geometry.coordinates,
          zoom,
        });
      });
  });
}

// =====================
// DOM MARKERS (ICONS FIX)
// =====================
function renderMarkers(reports) {
  // cleanup
  window.__markers.forEach((m) => m.remove());
  window.__markers = [];

  reports.forEach((r) => {
    const cat = r.category || "other";
    const icon = categoryIcon[cat] || "📍";

    const el = document.createElement("div");

    el.style.width = "32px";
    el.style.height = "32px";
    el.style.display = "flex";
    el.style.alignItems = "center";
    el.style.justifyContent = "center";
    el.style.fontSize = "22px";
    el.style.cursor = "pointer";
    el.style.userSelect = "none";

    el.textContent = icon;

    // ❌ NE zoomoljon, NE focusReport
    el.addEventListener("click", (e) => {
      e.stopPropagation();

      new maplibregl.Popup()
        .setLngLat([r.longitude, r.latitude])
        .setHTML(
          `
          <strong>${r.title}</strong><br/>
          ${r.description}<br/>
          <b>${categoryIcon[r.category] || "📍"} ${
            categoryLabel[r.category] || r.category
          }</b><br/>
          <b>${r.status}</b><br/>
          ${r.tracking_code}
        `,
        )
        .addTo(publicMap);
    });

    const marker = new maplibregl.Marker({
      element: el,
      anchor: "center",
    })
      .setLngLat([r.longitude, r.latitude])
      .addTo(publicMap);

    window.__markers.push(marker);
  });
}

// =====================
// GEO UPDATE (for clusters)
// =====================
function updateClusterData(reports) {
  const source = publicMap.getSource("reports");
  if (!source) return;

  source.setData({
    type: "FeatureCollection",
    features: reports.map((r) => ({
      type: "Feature",
      properties: {
        title: r.title,
        description: r.description,
        category: r.category,
        status: r.status,
        tracking_code: r.tracking_code,
      },
      geometry: {
        type: "Point",
        coordinates: [r.longitude, r.latitude],
      },
    })),
  });
}

// =====================
// FOCUS
// =====================
function focusReport(code) {
  const r = reportData.find((x) => x.tracking_code === code);
  if (!r) return;

  publicMap.flyTo({
    center: [r.longitude, r.latitude],
    zoom: 15,
  });

  new maplibregl.Popup()
    .setLngLat([r.longitude, r.latitude])
    .setHTML(
      `
      <strong>${r.title}</strong><br/>
      ${r.description}<br/>
      <b>${categoryIcon[r.category] || "📍"} ${
        categoryLabel[r.category] || r.category
      }</b><br/>
      <b>${r.status}</b><br/>
      ${r.tracking_code}
    `,
    )
    .addTo(publicMap);
}

// =====================
// LIST
// =====================
function renderList(data) {
  document.getElementById("list").innerHTML = data
    .map((r) => {
      const cat = r.category || "other";

      return `
      <div class="card" onclick="focusReport('${r.tracking_code}')">
        <b>${r.title}</b>
        <div class="small">${r.description}</div>

        <div class="badge ${r.status}">${r.status}</div>

        <div class="small">
          ${categoryIcon[cat] || "📍"} ${categoryLabel[cat] || cat}
        </div>

        <div class="small">${r.tracking_code}</div>
      </div>
    `;
    })
    .join("");
}

function renderStatusSummary(reportData) {
  const stats = calculateStats(reportData);

  const container = document.getElementById("status");
  if (!container) return;

  container.innerHTML = `
    <div class="status-summary ${activeStatusFilter ? "filtered" : ""}">

      <!-- ALL -->
      <div class="status-box all ${!activeStatusFilter ? "active" : ""}"
           onclick="setStatusFilter('all')">
        <strong>${reportData.length}</strong>
        <span>Összes</span>
      </div>

      <!-- OPEN -->
      <div class="status-box open ${activeStatusFilter === "open" ? "active" : ""}"
           onclick="setStatusFilter('open')">
        <strong>${stats.open}</strong>
        <span>Nyitott</span>
      </div>

      <!-- IN PROGRESS -->
      <div class="status-box in_progress ${activeStatusFilter === "in_progress" ? "active" : ""}"
           onclick="setStatusFilter('in_progress')">
        <strong>${stats.in_progress}</strong>
        <span>Folyamatban</span>
      </div>

      <!-- RESOLVED -->
      <div class="status-box resolved ${activeStatusFilter === "resolved" ? "active" : ""}"
           onclick="setStatusFilter('resolved')">
        <strong>${stats.resolved}</strong>
        <span>Megoldott</span>
      </div>

      <!-- INVALID -->
      <div class="status-box invalid ${activeStatusFilter === "invalid" ? "active" : ""}"
           onclick="setStatusFilter('invalid')">
        <strong>${stats.invalid}</strong>
        <span>Érvénytelen</span>
      </div>

    </div>
  `;
}

// =====================
// PICKER MAP
// =====================
function initPickerMap() {
  if (pickerMap) return;

  pickerMap = new maplibregl.Map({
    container: "picker-map",
    style:
      "https://tiles.basemaps.cartocdn.com/gl/dark-matter-gl-style/style.json",
    center: [19.04, 47.49],
    zoom: 6,
  });

  pickerMap.on("click", (e) => {
    selectedLng = e.lngLat.lng;
    selectedLat = e.lngLat.lat;

    if (pickerMarker) pickerMarker.remove();

    pickerMarker = new maplibregl.Marker({ color: "#ef4444" })
      .setLngLat([selectedLng, selectedLat])
      .addTo(pickerMap);
  });
}

// =====================
// MODAL
// =====================
function openModal() {
  const modal = document.getElementById("modal");
  modal.classList.add("show");
  initPickerMap();
}

function closeModal() {
  const modal = document.getElementById("modal");
  modal.classList.remove("show");
}

// =====================
// FILTER LOGIC
// =====================
function setStatusFilter(status) {
  if (status === "all") {
    activeStatusFilter = null;
  } else {
    activeStatusFilter = activeStatusFilter === status ? null : status;
  }

  applyFilters();
}

function applyFilters() {
  let filtered = reportData;

  if (activeStatusFilter) {
    filtered = filtered.filter((r) => r.status === activeStatusFilter);
  }

  renderList(filtered);
  renderMarkers(filtered);
  renderStatusSummary(reportData);
}

// =====================
// INIT
// =====================
function main() {
  initPublicMap();

  document.getElementById("new-report-btn").onclick = openModal;

  reportTitle.addEventListener("input", () => {
    document.getElementById("title-counter").textContent =
      `${reportTitle.value.length} / ${TITLE_LIMIT}`;
  });

  reportDescription.addEventListener("input", () => {
    document.getElementById("desc-counter").textContent =
      `${reportDescription.value.length} / ${DESC_LIMIT}`;
  });

  document.getElementById("report-form").addEventListener("submit", (e) => {
    e.preventDefault();
    createReport();
  });
}

main();
