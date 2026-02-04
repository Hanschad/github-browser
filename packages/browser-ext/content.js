// GitHub Browser - Content Script
// 在 GitHub 页面注入 "Open in IDE" 按钮

(function() {
  'use strict';

  // 检测页面类型
  function detectPageType() {
    const path = window.location.pathname;
    
    if (path.match(/^\/[^\/]+\/[^\/]+\/pull\/\d+/)) {
      return 'pull_request';
    } else if (path.match(/^\/[^\/]+\/[^\/]+\/blob\//)) {
      return 'file';
    } else if (path.match(/^\/[^\/]+\/[^\/]+\/tree\//)) {
      return 'directory';
    } else if (path.match(/^\/[^\/]+\/[^\/]+\/?$/)) {
      return 'repository';
    }
    
    return 'unknown';
  }

  // 创建按钮
  function createOpenButton() {
    const button = document.createElement('button');
    button.className = 'btn btn-sm github-browser-btn';
    button.innerHTML = `
      <svg class="octicon" width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
        <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
      </svg>
      Open in IDE
    `;
    button.title = 'Open in local IDE with GitHub Browser';
    
    button.addEventListener('click', async (e) => {
      e.preventDefault();
      e.stopPropagation();
      
      const url = window.location.href;
      await openInIDE(url);
    });
    
    return button;
  }

  // 打开 IDE
  async function openInIDE(url) {
    try {
      // 显示加载状态
      showNotification('Opening in IDE...', 'info');

      // 通过 background script 发送请求（避免 Firefox CORS 限制）
      const response = await chrome.runtime.sendMessage({
        action: 'openInIDE',
        url: url
      });

      if (response && response.success) {
        showNotification('✅ Opened successfully!', 'success');
      } else if (response && response.error) {
        throw new Error(response.error);
      } else {
        throw new Error('Unknown error');
      }

    } catch (error) {
      console.error('GitHub Browser error:', error);
      
      if (error.message.includes('fetch') || error.message.includes('network') || error.message.includes('NetworkError')) {
        showNotification(
          '❌ Cannot connect to GitHub Browser service. Make sure it is running.',
          'error'
        );
      } else {
        showNotification(`❌ Error: ${error.message}`, 'error');
      }
    }
  }

  // 显示通知
  function showNotification(message, type = 'info') {
    // 移除旧通知
    const oldNotif = document.querySelector('.github-browser-notification');
    if (oldNotif) {
      oldNotif.remove();
    }

    // 创建新通知
    const notif = document.createElement('div');
    notif.className = `github-browser-notification github-browser-notification-${type}`;
    notif.textContent = message;
    document.body.appendChild(notif);

    // 3 秒后自动移除
    setTimeout(() => {
      notif.style.opacity = '0';
      setTimeout(() => notif.remove(), 300);
    }, 3000);
  }

  // 注入按钮
  function injectButton() {
    const pageType = detectPageType();
    
    if (pageType === 'unknown') {
      return;
    }

    // 检查是否已经注入
    if (document.querySelector('.github-browser-btn')) {
      return;
    }

    const button = createOpenButton();
    let inserted = false;

    if (pageType === 'pull_request') {
      // PR 页面：放在 Edit/Code 按钮旁边
      
      // 查找右上角的 Code 下拉按钮（带有 aria-haspopup 的按钮，且可见）
      const allButtons = document.querySelectorAll('button');
      for (const btn of allButtons) {
        const text = btn.textContent.trim();
        // 匹配 Code 按钮：文本以 Code 开头，有下拉菜单属性，是小尺寸按钮
        if ((text === 'Code' || text.startsWith('Code')) && 
            btn.getAttribute('aria-haspopup') === 'true' &&
            btn.getAttribute('data-size') === 'small') {
          const parent = btn.parentElement;
          // 检查父元素是否可见（offsetHeight > 0）
          if (parent && parent.offsetHeight > 0) {
            btn.before(button);
            inserted = true;
            break;
          }
        }
      }
      
      // 备选：查找可见的 Edit 按钮
      if (!inserted) {
        for (const btn of allButtons) {
          const text = btn.textContent.trim();
          if (text === 'Edit' && btn.getAttribute('data-size') === 'small') {
            const parent = btn.parentElement;
            if (parent && parent.offsetHeight > 0) {
              btn.after(button);
              inserted = true;
              break;
            }
          }
        }
      }
      
      // 备选：通过选择器查找
      if (!inserted) {
        const targets = [
          '.gh-header-actions',
          '[data-testid="issue-header-actions"]'
        ];
        for (const sel of targets) {
          const target = document.querySelector(sel);
          if (target && target.offsetHeight > 0) {
            target.prepend(button);
            inserted = true;
            break;
          }
        }
      }
    } else if (pageType === 'file') {
      // 文件页面：放在文件操作栏 (Raw/Blame 按钮旁边)
      const targets = [
        '.react-blob-header-edit-and-raw-actions',
        '[data-testid="raw-button"]',
        '.Box-header .d-flex.gap-2',
        '.Box-header button[data-hotkey]'
      ];
      for (const sel of targets) {
        const target = document.querySelector(sel);
        if (target) {
          const parent = target.closest('.d-flex, .BtnGroup') || target.parentElement;
          parent.prepend(button);
          inserted = true;
          break;
        }
      }
    } else if (pageType === 'repository' || pageType === 'directory') {
      // 仓库主页/目录页：放在 Code 按钮旁边
      
      // 方法1：查找包含 "Code" 文字的按钮（最可靠）
      const allButtons = document.querySelectorAll('button, summary');
      for (const btn of allButtons) {
        const text = btn.textContent.trim();
        // 匹配 "Code" 按钮（通常是绿色的下拉按钮）
        if ((text === 'Code' || text.startsWith('Code')) && 
            (btn.classList.contains('btn-primary') || btn.closest('.btn-primary') || 
             btn.querySelector('.octicon-code') || btn.closest('[data-component="IconButton"]'))) {
          const parent = btn.closest('.d-flex, .gap-2, .BtnGroup') || btn.parentElement;
          if (parent) {
            // 在 Code 按钮前面插入
            btn.before(button);
            inserted = true;
            break;
          }
        }
      }
      
      // 方法2：通过选择器查找
      if (!inserted) {
        const targets = [
          // React 版本的 Code 按钮容器
          '[class*="react-code-view"] .d-flex.gap-2',
          // get-repo 组件
          'get-repo',
          '#repo-content-pjax-container .d-flex.gap-2',
          // 文件导航区域
          '.file-navigation .d-flex.gap-2',
          '.Box-header .d-flex.gap-2'
        ];
        for (const sel of targets) {
          const target = document.querySelector(sel);
          if (target) {
            target.prepend(button);
            inserted = true;
            break;
          }
        }
      }
      
      // 方法3：查找 Edit 按钮，在其旁边插入
      if (!inserted) {
        const editBtn = Array.from(document.querySelectorAll('button, a')).find(
          el => el.textContent.trim() === 'Edit' || el.textContent.trim() === 'Go to file'
        );
        if (editBtn) {
          const parent = editBtn.closest('.d-flex, .gap-2') || editBtn.parentElement;
          if (parent) {
            parent.appendChild(button);
            inserted = true;
          }
        }
      }
    }

    // 最终备选：固定在右上角
    if (!inserted) {
      button.classList.add('github-browser-btn-fixed');
      document.body.appendChild(button);
    }
  }

  // 监听 URL 变化（GitHub 是 SPA）
  let lastUrl = location.href;
  new MutationObserver(() => {
    const url = location.href;
    if (url !== lastUrl) {
      lastUrl = url;
      setTimeout(injectButton, 500);
    }
  }).observe(document, { subtree: true, childList: true });

  // 初始注入
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
      setTimeout(injectButton, 1000);
    });
  } else {
    setTimeout(injectButton, 1000);
  }

  // 添加键盘快捷键：Shift+O
  document.addEventListener('keydown', (e) => {
    if (e.shiftKey && e.key === 'O') {
      e.preventDefault();
      openInIDE(window.location.href);
    }
  });

})();
