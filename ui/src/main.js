const { invoke } = window.__TAURI__.core;
let searchInput;

async function performSearch(query) {
  if (!query.trim()) return;

  const url = `http://localhost:8080/?q=${encodeURIComponent(query)}`;
  await invoke('open_url', { url });

  searchInput.value = '';
  setTimeout(() => invoke('hide_window'), 100);
}

function handleKeydown(e) {
  if (e.key === "Escape") {
    invoke('hide_window');
  }
}

window.addEventListener("DOMContentLoaded", () => {
  searchInput = document.querySelector("#search-input");

  document.querySelector("#search-form")?.addEventListener("submit", (e) => {
    e.preventDefault();
    performSearch(searchInput.value);
  });

  document.addEventListener("keydown", handleKeydown);
  searchInput?.focus();
});
