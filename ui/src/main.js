const { invoke } = window.__TAURI__.core;
const { WebviewWindow } = window.__TAURI__.webviewWindow;

let searchInputEl;

async function performSearch(query) {
  if (!query.trim()) return;
  
  // Open the browser with the search
  const url = `http://localhost:8080/?q=${encodeURIComponent(query)}`;
  await invoke('open_url', { url });
  
  // Clear the input and hide the window after search
  searchInputEl.value = '';
  const currentWindow = WebviewWindow.getCurrentWebviewWindow();
  
  // Small delay for better UX, then hide
  setTimeout(async () => {
    await currentWindow.hide();
  }, 100);
}

window.addEventListener("DOMContentLoaded", () => {
  searchInputEl = document.querySelector("#search-input");
  
  document.querySelector("#search-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    await performSearch(searchInputEl.value);
  });
  
  // ESC key to hide window
  document.addEventListener("keydown", async (e) => {
    if (e.key === "Escape") {
      const currentWindow = WebviewWindow.getCurrentWebviewWindow();
      await currentWindow.hide();
    }
  });
  
  // Focus the input when window is shown
  searchInputEl.focus();
});