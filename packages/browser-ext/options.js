// GitHub Browser - Options Script

// 加载设置
async function loadSettings() {
  const config = await chrome.storage.sync.get({
    serviceUrl: 'http://localhost:9527',
    ide: 'code'
  });

  document.getElementById('service-url').value = config.serviceUrl;
  document.getElementById('ide').value = config.ide;
}

// 保存设置
async function saveSettings(e) {
  e.preventDefault();

  const serviceUrl = document.getElementById('service-url').value;
  const ide = document.getElementById('ide').value;

  try {
    // 测试连接
    const response = await fetch(`${serviceUrl}/health`, {
      method: 'GET',
      signal: AbortSignal.timeout(2000)
    });

    if (!response.ok) {
      throw new Error('Service returned error');
    }

    // 保存配置
    await chrome.storage.sync.set({
      serviceUrl: serviceUrl,
      ide: ide
    });

    // 显示成功消息
    showStatus('Settings saved successfully!', 'success');

  } catch (error) {
    showStatus(
      'Warning: Settings saved, but cannot connect to service. Make sure it is running.',
      'error'
    );
  }
}

// 显示状态消息
function showStatus(message, type) {
  const statusEl = document.getElementById('status');
  statusEl.textContent = message;
  statusEl.className = `status ${type}`;

  setTimeout(() => {
    statusEl.className = 'status';
  }, 3000);
}

// 初始化
document.addEventListener('DOMContentLoaded', () => {
  loadSettings();
  document.getElementById('settings-form').addEventListener('submit', saveSettings);
});
