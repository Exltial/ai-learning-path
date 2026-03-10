import React, { useState, useEffect } from 'react';
import { Download, X } from 'lucide-react';
import './PWAInstallPrompt.css';

interface BeforeInstallPromptEvent extends Event {
  prompt: () => void;
  userChoice: Promise<{ outcome: 'accepted' | 'dismissed' }>;
}

const PWAInstallPrompt: React.FC = () => {
  const [deferredPrompt, setDeferredPrompt] = useState<BeforeInstallPromptEvent | null>(null);
  const [showPrompt, setShowPrompt] = useState(false);

  useEffect(() => {
    // 检查是否已经安装过
    const hasDismissed = localStorage.getItem('pwa-install-dismissed');
    if (hasDismissed) {
      return;
    }

    const handleBeforeInstallPrompt = (e: Event) => {
      e.preventDefault();
      const promptEvent = e as BeforeInstallPromptEvent;
      setDeferredPrompt(promptEvent);
      
      // 延迟显示提示，让用户先体验应用
      setTimeout(() => {
        setShowPrompt(true);
      }, 3000);
    };

    window.addEventListener('beforeinstallprompt', handleBeforeInstallPrompt);

    return () => {
      window.removeEventListener('beforeinstallprompt', handleBeforeInstallPrompt);
    };
  }, []);

  const handleInstall = async () => {
    if (!deferredPrompt) return;

    try {
      deferredPrompt.prompt();
      const { outcome } = await deferredPrompt.userChoice;
      
      if (outcome === 'accepted') {
        console.log('用户同意安装 PWA');
      } else {
        console.log('用户拒绝安装 PWA');
        localStorage.setItem('pwa-install-dismissed', 'true');
      }
      
      setShowPrompt(false);
      setDeferredPrompt(null);
    } catch (error) {
      console.error('安装 PWA 失败:', error);
    }
  };

  const handleDismiss = () => {
    setShowPrompt(false);
    localStorage.setItem('pwa-install-dismissed', 'true');
  };

  // 检查是否在 iOS 上
  const isIOS = /iPad|iPhone|iPod/.test(navigator.userAgent) && !(window as any).MSStream;
  
  // 检查是否在 standalone 模式
  const isStandalone = window.matchMedia('(display-mode: standalone)').matches;

  // 如果已经安装或用户已关闭，不显示
  if (isStandalone || !showPrompt) {
    return null;
  }

  return (
    <div className="pwa-install-prompt">
      <div className="pwa-install-prompt-content">
        <div className="pwa-install-prompt-icon">
          <Download size={32} className="text-primary-600" />
        </div>
        <div className="pwa-install-prompt-text">
          <h3 className="pwa-install-prompt-title">
            安装 AI 学习平台
          </h3>
          <p className="pwa-install-prompt-description">
            添加到主屏幕，离线也能访问，享受更好的学习体验
          </p>
        </div>
        <button
          className="pwa-install-prompt-close"
          onClick={handleDismiss}
          aria-label="关闭安装提示"
        >
          <X size={20} />
        </button>
      </div>
      <div className="pwa-install-prompt-actions">
        {!isIOS ? (
          <button className="pwa-install-prompt-btn" onClick={handleInstall}>
            立即安装
          </button>
        ) : (
          <div className="pwa-install-prompt-ios-instructions">
            <p>在 Safari 中：</p>
            <ol>
              <li>点击底部 <strong>分享</strong> 按钮</li>
              <li>选择 <strong>添加到主屏幕</strong></li>
            </ol>
          </div>
        )}
        <button className="pwa-install-prompt-dismiss" onClick={handleDismiss}>
          暂时不用
        </button>
      </div>
    </div>
  );
};

export default PWAInstallPrompt;
