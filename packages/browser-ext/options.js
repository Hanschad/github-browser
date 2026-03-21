async function loadSettings() {
  const stored = await chrome.storage.sync.get({
    serviceUrl: "http://localhost:9527",
    ide: "code",
    pathMappings: []
  });

  let config = stored;
  try {
    const backendResponse = await chrome.runtime.sendMessage({
      action: "getBackendConfig",
      serviceUrl: stored.serviceUrl
    });

    if (backendResponse && backendResponse.success && backendResponse.result.config) {
      config = {
        ...stored,
        ide: backendResponse.result.config.defaultIDE || stored.ide,
        pathMappings: backendResponse.result.config.pathMappings || stored.pathMappings
      };

      await chrome.storage.sync.set({
        serviceUrl: config.serviceUrl,
        ide: config.ide,
        pathMappings: config.pathMappings
      });
    }
  } catch (error) {
    config = stored;
  }

  document.getElementById("service-url").value = config.serviceUrl;
  document.getElementById("ide").value = config.ide;

  const container = document.getElementById("path-mappings");
  container.innerHTML = "";
  if (config.pathMappings.length === 0) {
    addMappingRow("", "");
  } else {
    config.pathMappings.forEach((mapping) => addMappingRow(mapping.pattern, mapping.localPath));
  }
}

function addMappingRow(pattern = "", localPath = "") {
  const container = document.getElementById("path-mappings");
  const row = document.createElement("div");
  row.className = "mapping-row";
  row.innerHTML = `
    <input type="text" class="mapping-pattern" placeholder="microsoft or */repo" value="${pattern}">
    <input type="text" class="mapping-path" placeholder="~/projects/microsoft" value="${localPath}">
    <button type="button" class="remove-mapping">×</button>
  `;
  row.querySelector(".remove-mapping").addEventListener("click", () => {
    row.remove();
  });
  container.appendChild(row);
}

function getPathMappings() {
  const rows = document.querySelectorAll(".mapping-row");
  const mappings = [];
  rows.forEach((row) => {
    const pattern = row.querySelector(".mapping-pattern").value.trim();
    const localPath = row.querySelector(".mapping-path").value.trim();
    if (pattern && localPath) {
      mappings.push({ pattern, localPath });
    }
  });
  return mappings;
}

async function saveSettings(e) {
  e.preventDefault();

  const serviceUrl = document.getElementById("service-url").value.trim();
  const ide = document.getElementById("ide").value;
  const pathMappings = getPathMappings();

  await chrome.storage.sync.set({
    serviceUrl,
    ide,
    pathMappings
  });

  try {
    const configResponse = await chrome.runtime.sendMessage({
      action: "getBackendConfig",
      serviceUrl
    });

    if (!configResponse || !configResponse.success) {
      throw new Error((configResponse && configResponse.error) || "Cannot load backend config");
    }

    const nextConfig = {
      ...configResponse.result.config,
      defaultIDE: ide,
      pathMappings
    };

    const updateResponse = await chrome.runtime.sendMessage({
      action: "updateBackendConfig",
      serviceUrl,
      config: nextConfig
    });

    if (!updateResponse || !updateResponse.success) {
      throw new Error((updateResponse && updateResponse.error) || "Cannot update backend config");
    }

    showStatus("Settings saved and synced successfully.", "success");
  } catch (error) {
    showStatus(
      `Settings saved locally, but backend sync failed: ${error.message}`,
      "error"
    );
  }
}

function showStatus(message, type) {
  const statusEl = document.getElementById("status");
  statusEl.textContent = message;
  statusEl.className = `status ${type}`;

  setTimeout(() => {
    statusEl.className = "status";
  }, 3000);
}

document.addEventListener("DOMContentLoaded", () => {
  loadSettings();
  document.getElementById("settings-form").addEventListener("submit", saveSettings);
  document.getElementById("add-mapping").addEventListener("click", () => addMappingRow());
});
