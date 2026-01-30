// GitHub Browser - Background Script

// 监听来自 content script 的消息
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.action === 'openInIDE') {
    handleOpenInIDE(request.url)
      .then(result => sendResponse({ success: true, result }))
      .catch(error => sendResponse({ success: false, error: error.message }));
    return true; // 保持消息通道开放
  }
});

// 处理打开 IDE 请求
async function handleOpenInIDE(url) {
  const config = await chrome.storage.sync.get({
    serviceUrl: 'http://localhost:9527',
    ide: 'code'
  });

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

  return await response.json();
}
