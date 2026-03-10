import React, { useState, useEffect } from 'react';
import { Home, BookOpen, BarChart3, Trophy, User, Code } from 'lucide-react';
import { useLocation, useNavigate } from 'react-router-dom';
import './MobileNav.css';

interface NavItem {
  path: string;
  icon: React.ReactNode;
  label: string;
}

const navItems: NavItem[] = [
  { path: '/', icon: <Home size={24} />, label: '首页' },
  { path: '/courses', icon: <BookOpen size={24} />, label: '课程' },
  { path: '/practice', icon: <Code size={24} />, label: '练习' },
  { path: '/progress', icon: <BarChart3 size={24} />, label: '进度' },
  { path: '/leaderboard', icon: <Trophy size={24} />, label: '排行' },
];

const MobileNav: React.FC = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState(location.pathname);
  const [isVisible, setIsVisible] = useState(true);
  const [lastScrollY, setLastScrollY] = useState(0);

  // 监听路由变化
  useEffect(() => {
    setActiveTab(location.pathname);
  }, [location.pathname]);

  // 监听滚动，控制导航栏显示/隐藏
  useEffect(() => {
    const handleScroll = () => {
      const currentScrollY = window.scrollY;
      
      // 向下滚动时隐藏，向上滚动时显示
      if (currentScrollY > lastScrollY && currentScrollY > 100) {
        setIsVisible(false);
      } else {
        setIsVisible(true);
      }
      
      setLastScrollY(currentScrollY);
    };

    window.addEventListener('scroll', handleScroll, { passive: true });
    return () => window.removeEventListener('scroll', handleScroll);
  }, [lastScrollY]);

  const handleNavClick = (path: string) => {
    if (path !== activeTab) {
      navigate(path);
    }
  };

  // 某些页面不显示底部导航
  const hiddenPaths = ['/login', '/register', '/settings'];
  if (hiddenPaths.some(path => location.pathname.startsWith(path))) {
    return null;
  }

  return (
    <nav 
      className={`mobile-nav ${isVisible ? 'mobile-nav-visible' : 'mobile-nav-hidden'}`}
      role="navigation"
      aria-label="底部导航"
    >
      <div className="mobile-nav-container">
        {navItems.map((item) => {
          const isActive = activeTab === item.path || 
            (item.path !== '/' && activeTab.startsWith(item.path));
          
          return (
            <button
              key={item.path}
              className={`mobile-nav-item ${isActive ? 'active' : ''}`}
              onClick={() => handleNavClick(item.path)}
              aria-label={item.label}
              aria-current={isActive ? 'page' : undefined}
            >
              <div className="mobile-nav-icon">
                {item.icon}
              </div>
              <span className="mobile-nav-label">{item.label}</span>
              {isActive && (
                <span className="mobile-nav-indicator" aria-hidden="true" />
              )}
            </button>
          );
        })}
      </div>
      
      {/* 安全区域适配 iOS 底部横条 */}
      <div className="mobile-nav-safe-area" aria-hidden="true" />
    </nav>
  );
};

export default MobileNav;
