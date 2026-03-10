import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './index.css'
import './styles/mobile.css'

// 注册 PWA Service Worker
if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('/sw.js', {
      scope: '/'
    }).then(registration => {
      console.log('SW registered:', registration.scope);
      
      // 检查更新
      registration.addEventListener('updatefound', () => {
        const newWorker = registration.installing;
        if (newWorker) {
          newWorker.addEventListener('statechange', () => {
            if (newWorker.state === 'installed' && navigator.serviceWorker.controller) {
              // 有新版本可用
              console.log('New content available, please refresh.');
              if (confirm('有新版本可用，是否刷新？')) {
                window.location.reload();
              }
            }
          });
        }
      });
    }).catch(error => {
      console.log('SW registration failed:', error);
    });
  });
}

// 请求通知权限
if ('Notification' in window && Notification.permission === 'default') {
  // 在用户交互后请求权限
  document.addEventListener('click', () => {
    Notification.requestPermission();
  }, { once: true });
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
