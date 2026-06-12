export function setStatus(message, type = "info") {
  const container = document.getElementById("toast-container");

  const icons = {
    success: "✓",
    error: "✕",
    info: "ℹ",
  };

  const toast = document.createElement("div");
  toast.className = `toast ${type}`;

  toast.innerHTML = `
    <div class="toast-icon">${icons[type] || "ℹ"}</div>
    <div class="toast-message">${message}</div>
  `;

  container.appendChild(toast);
  setTimeout(() => toast.remove(), 4000);
}
