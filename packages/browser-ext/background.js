const DEFAULT_SERVICE_URL = "http://localhost:9527";
const NATIVE_HOST_NAME = "com.github.browser";

chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  handleRequest(request)
    .then((result) => sendResponse({ success: true, result }))
    .catch((error) => sendResponse({ success: false, error: error.message }));

  return true;
});

async function handleRequest(request) {
  switch (request.action) {
    case "openInIDE":
      return openInIDE(request.url);
    case "checkBackend":
      return checkBackend(request.serviceUrl);
    case "getBackendConfig":
      return getBackendConfig(request.serviceUrl);
    case "updateBackendConfig":
      return updateBackendConfig(request.config, request.serviceUrl);
    default:
      throw new Error(`Unknown action: ${request.action}`);
  }
}

async function openInIDE(url) {
  const config = await chrome.storage.sync.get({
    serviceUrl: DEFAULT_SERVICE_URL,
    ide: "code"
  });

  return callBackend({
    action: "open",
    url,
    ide: config.ide
  }, config.serviceUrl);
}

async function checkBackend(serviceUrl) {
  try {
    const response = await callNativeHost({ action: "health" });
    return {
      ...response,
      transport: "native"
    };
  } catch (nativeError) {
    const httpResponse = await callHttpBackend({ action: "health" }, serviceUrl, nativeError);
    return {
      ...httpResponse,
      transport: "http"
    };
  }
}

async function getBackendConfig(serviceUrl) {
  const response = await callBackend({ action: "getConfig" }, serviceUrl);
  return {
    ...response,
    config: response.config || null
  };
}

async function updateBackendConfig(config, serviceUrl) {
  if (!config) {
    throw new Error("Config is required");
  }

  return callBackend({
    action: "updateConfig",
    config
  }, serviceUrl);
}

async function callBackend(message, serviceUrl) {
  try {
    const nativeResponse = await callNativeHost(message);
    return {
      ...nativeResponse,
      transport: "native"
    };
  } catch (nativeError) {
    if (nativeError.backendKind !== "transport") {
      throw nativeError;
    }

    const httpResponse = await callHttpBackend(message, serviceUrl, nativeError);
    return {
      ...httpResponse,
      transport: "http"
    };
  }
}

function callNativeHost(message) {
  return new Promise((resolve, reject) => {
    if (!chrome.runtime.connectNative) {
      reject(createBackendError("Native messaging is not supported in this browser", "transport"));
      return;
    }

    let settled = false;
    let port;

    try {
      port = chrome.runtime.connectNative(NATIVE_HOST_NAME);
    } catch (error) {
      reject(createBackendError(error.message, "transport"));
      return;
    }

    function fail(message, kind) {
      if (settled) {
        return;
      }

      settled = true;
      try {
        port.disconnect();
      } catch (error) {}
      reject(createBackendError(message, kind));
    }

    port.onMessage.addListener((response) => {
      if (settled) {
        return;
      }

      settled = true;
      try {
        port.disconnect();
      } catch (error) {}

      if (!response) {
        reject(createBackendError("Native host returned an empty response", "transport"));
        return;
      }

      if (response.status && response.status !== "ok") {
        reject(createBackendError(response.message || "Native host request failed", "request"));
        return;
      }

      resolve(response);
    });

    port.onDisconnect.addListener(() => {
      if (settled) {
        return;
      }

      const message = chrome.runtime.lastError
        ? chrome.runtime.lastError.message
        : "Native host disconnected";
      fail(message, "transport");
    });

    try {
      port.postMessage(message);
    } catch (error) {
      fail(error.message, "transport");
    }
  });
}

async function callHttpBackend(message, serviceUrl, nativeError) {
  const baseUrl = serviceUrl || DEFAULT_SERVICE_URL;

  try {
    switch (message.action) {
      case "health":
        return await httpGetJSON(`${baseUrl}/health`);
      case "open":
        return await httpJSON(`${baseUrl}/open`, "POST", {
          url: message.url,
          ide: message.ide
        });
      case "getConfig": {
        const config = await httpGetJSON(`${baseUrl}/config`);
        return {
          status: "ok",
          config
        };
      }
      case "updateConfig":
        return await httpJSON(`${baseUrl}/config`, "PUT", message.config);
      default:
        throw new Error(`HTTP fallback does not support action: ${message.action}`);
    }
  } catch (error) {
    const prefix = nativeError
      ? `Native host unavailable (${nativeError.message}). `
      : "";
    throw new Error(`${prefix}${error.message}`);
  }
}

async function httpGetJSON(url) {
  const response = await fetch(url, {
    method: "GET",
    signal: AbortSignal.timeout(2000)
  });
  return parseJSONResponse(response);
}

async function httpJSON(url, method, body) {
  const response = await fetch(url, {
    method,
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(body),
    signal: AbortSignal.timeout(5000)
  });
  return parseJSONResponse(response);
}

async function parseJSONResponse(response) {
  let payload = null;

  try {
    payload = await response.json();
  } catch (error) {}

  if (!response.ok) {
    throw new Error((payload && payload.message) || (payload && payload.error) || "Backend request failed");
  }

  return payload || { status: "ok" };
}

function createBackendError(message, kind) {
  const error = new Error(message);
  error.backendKind = kind;
  return error;
}
