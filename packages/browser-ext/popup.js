// GitHub Browser - Popup Script

let config = {
  serviceUrl: 'http://localhost:9527',
  ide: 'code'
};

// 加载配置
async function loadConfig() {
  const stored = await chrome.storage.sync.get({
    serviceUrl: 'http://localhost:9527',
    ide: 'code'
  });
  config = stored;
}

// 检查服务状态
async function checkServiceStatus() {
  const statusEl = document.getElementById('status');
  const statusTextEl = document.getElementById('status-text');

  try {
    const response = await fetch(`${config.serviceUrl}/health`, {
      method: 'GET',
      signal: AbortSignal.timeout(2000)
    });

    if (response.ok) {
      const data = await response.json();
      statusEl.className = 'status status-ok';
      statusTextEl.textContent = `Service running (v${data.version})`;
    } else {
      throw new Error('Service returned error');
    }
  } catch (error) {
    statusEl.className = 'status status-error';
    statusTextEl.textContent = 'Service not running';
  }
}

// 打开当前页面
async function openCurrentPage() {
  const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
  
  if (!tab.url || !tab.url.includes('github.com')) {
    alert('Please open a GitHub page first');
    return;
  }

  await openInIDE(tab.url);
}

// 从剪贴板打开
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

// 打开 IDE
async function openInIDE(url) {
  try {
    const response = await fetch(`${config.serviceUrl}/open`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        url: url,
        ide: config.ide
      })
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || 'Failed to open repository');
    }

    const result = await response.json();
    
    // 显示成功消息
    const statusEl = document.getElementById('status');
    const statusTextEl = document.getElementById('status-text');
    statusEl.className = 'status status-ok';
    statusTextEl.textContent = '✅ Opened successfully!';

    // 3 秒后关闭 popup
    setTimeout(() => {
      window.close();
    }, 1500);

  } catch (error) {
    alert('Error: ' + error.message);
  }
}

// 打开设置
function openOptions() {
  chrome.runtime.openOptionsPage();
}

// 初始化
document.addEventListener('DOMContentLoaded', async () => {
  await loadConfig();
  await checkServiceStatus();

  document.getElementById('open-current').addEventListener('click', openCurrentPage);
  document.getElementById('open-from-clipboard').addEventListener('click', openFromClipboard);
  document.getElementById('open-options').addEventListener('click', openOptions);
});
