import React, { useState, useEffect } from 'react';
import { WifiOff, RefreshCw, Home } from 'lucide-react';
import './OfflinePage.css';

const OfflinePage: React.FC = () => {
  const [isOnline, setIsOnline] = useState(navigator.onLine);

  useEffect(() => {
    const handleOnline = () => setIsOnline(true);
    const handleOffline = () => setIsOnline(false);

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
    };
  }, []);

  const handleRefresh = () => {
    window.location.reload();
  };

  const handleGoHome = () => {
    window.location.href = '/';
  };

  if (isOnline) {
    return null;
  }

  return (
    <div className="offline-page">
      <div className="offline-page-content">
        <div className="offline-page-icon">
          <WifiOff size={64} />
        </div>
        
        <h1 className="offline-page-title">
          网络已断开
        </h1>
        
        <p className="offline-page-description">
          您已离线，但仍可访问已缓存的内容。
          <br />
          请检查网络连接后重试。
        </p>
        
        <div className="offline-page-actions">
          <button className="offline-page-btn primary" onClick={handleRefresh}>
            <RefreshCw size={20} />
            <span>重新连接</span>
          </button>
          
          <button className="offline-page-btn secondary" onClick={handleGoHome}>
            <Home size={20} />
            <span>返回首页</span>
          </button>
        </div>
        
        <div className="offline-page-tips">
          <h3>离线模式下您可以：</h3>
          <ul>
            <li>查看已缓存的课程内容</li>
            <li>浏览之前访问过的页面</li>
            <li>查看本地保存的学习进度</li>
          </ul>
        </div>
      </div>
    </div>
  );
};

export default OfflinePage;
