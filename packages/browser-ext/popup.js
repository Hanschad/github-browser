let config = {
  serviceUrl: "http://localhost:9527",
  ide: "code"
};

async function loadConfig() {
  const stored = await chrome.storage.sync.get({
    serviceUrl: "http://localhost:9527",
    ide: "code"
  });
  config = stored;
}

async function checkBackendStatus() {
  const statusEl = document.getElementById("status");
  const statusTextEl = document.getElementById("status-text");

  try {
    const response = await chrome.runtime.sendMessage({
      action: "checkBackend",
      serviceUrl: config.serviceUrl
    });

    if (response && response.success) {
      const backend = response.result.transport === "native" ? "Native host" : "HTTP service";
      statusEl.className = "status status-ok";
      statusTextEl.textContent = `${backend} ready (v${response.result.version || "unknown"})`;
      return;
    }

    statusEl.className = "status status-error";
    statusTextEl.textContent = (response && response.error) || "Backend unavailable";
  } catch (error) {
    statusEl.className = "status status-error";
    statusTextEl.textContent = error.message || "Backend unavailable";
  }
}

async function openCurrentPage() {
  const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });

  if (!tab.url || !tab.url.includes('github.com')) {
    alert('Please open a GitHub page first');
    return;
  }

  await openInIDE(tab.url);
}

async function openFromClipboard() {
  try {
    const text = await navigator.clipboard.readText();

    if (!text || !text.includes('github.com')) {
      alert('Clipboard does not contain a GitHub URL');
      return;
    }

    await openInIDE(text);
  } catch (error) {
    alert('Failed to read clipboard: ' + error.message);
  }
}

async function openInIDE(url) {
  const response = await chrome.runtime.sendMessage({
    action: "openInIDE",
    url
  });

  if (!response || !response.success) {
    throw new Error((response && response.error) || "Failed to open repository");
  }

  const statusEl = document.getElementById("status");
  const statusTextEl = document.getElementById("status-text");
  statusEl.className = "status status-ok";
  statusTextEl.textContent = "Opened successfully";

  setTimeout(() => {
    window.close();
  }, 1500);
}

function openOptions() {
  chrome.runtime.openOptionsPage();
}

document.addEventListener("DOMContentLoaded", async () => {
  await loadConfig();
  await checkBackendStatus();

  document.getElementById("open-current").addEventListener("click", async () => {
    try {
      await openCurrentPage();
    } catch (error) {
      alert("Error: " + error.message);
    }
  });

  document.getElementById("open-from-clipboard").addEventListener("click", async () => {
    try {
      await openFromClipboard();
    } catch (error) {
      alert("Error: " + error.message);
    }
  });

  document.getElementById("open-options").addEventListener("click", openOptions);
});
