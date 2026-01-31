// GitHub Browser - Options Script

// 加载设置
async function loadSettings() {
  const config = await chrome.storage.sync.get({
    serviceUrl: 'http://localhost:9527',
    ide: 'code',
    pathMappings: []
  });

  document.getElementById('service-url').value = config.serviceUrl;
  document.getElementById('ide').value = config.ide;
  
  // 加载路径映射
  const container = document.getElementById('path-mappings');
  container.innerHTML = '';
  if (config.pathMappings.length === 0) {
    addMappingRow('', '');
  } else {
    config.pathMappings.forEach(m => addMappingRow(m.pattern, m.localPath));
  }
}

// 添加路径映射行
function addMappingRow(pattern = '', localPath = '') {
  const container = document.getElementById('path-mappings');
  const row = document.createElement('div');
  row.className = 'mapping-row';
  row.innerHTML = `
    <input type="text" class="mapping-pattern" placeholder="microsoft or */repo" value="${pattern}">
    <input type="text" class="mapping-path" placeholder="~/projects/microsoft" value="${localPath}">
    <button type="button" class="remove-mapping">×</button>
  `;
  row.querySelector('.remove-mapping').addEventListener('click', () => {
    row.remove();
  });
  container.appendChild(row);
}

// 获取所有路径映射
function getPathMappings() {
  const rows = document.querySelectorAll('.mapping-row');
  const mappings = [];
  rows.forEach(row => {
    const pattern = row.querySelector('.mapping-pattern').value.trim();
    const localPath = row.querySelector('.mapping-path').value.trim();
    if (pattern && localPath) {
      mappings.push({ pattern, localPath });
    }
  });
  return mappings;
}

// 保存设置
async function saveSettings(e) {
  e.preventDefault();

  const serviceUrl = document.getElementById('service-url').value;
  const ide = document.getElementById('ide').value;
  const pathMappings = getPathMappings();

  try {
    // 测试连接
    const response = await fetch(`${serviceUrl}/health`, {
      method: 'GET',
      signal: AbortSignal.timeout(2000)
    });

    if (!response.ok) {
      throw new Error('Service returned error');
    }

    // 保存配置到浏览器扩展
    await chrome.storage.sync.set({
      serviceUrl: serviceUrl,
      ide: ide,
      pathMappings: pathMappings
    });

    // 同步路径映射配置到服务端
    if (pathMappings.length > 0) {
      try {
        // 先获取当前服务端配置
        const configRes = await fetch(`${serviceUrl}/config`);
        const currentConfig = await configRes.json();
        
        // 更新 pathMappings
        await fetch(`${serviceUrl}/config`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ ...currentConfig, pathMappings: pathMappings })
        });
      } catch (e) {
        // 忽略同步失败
      }
    }

    // 显示成功消息
    showStatus('Settings saved successfully!', 'success');

  } catch (error) {
    // 即使服务不可用也保存本地配置
    await chrome.storage.sync.set({
      serviceUrl: serviceUrl,
      ide: ide,
      pathMappings: pathMappings
    });
    
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
  document.getElementById('add-mapping').addEventListener('click', () => addMappingRow());
});
