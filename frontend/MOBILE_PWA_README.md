# 移动端优化与 PWA 实现文档

## 📱 实现功能

### 1. 响应式布局优化

- **断点设计**
  - 手机：< 768px
  - 平板：768px - 1023px
  - 桌面：≥ 1024px

- **移动端优化**
  - 触摸目标最小尺寸 48x48px
  - 字体大小适配（防止 iOS 自动缩放）
  - 安全区域适配（iOS 刘海屏/底部横条）
  - 横屏模式优化

- **文件位置**: `src/styles/mobile.css`

### 2. PWA (Progressive Web App)

- **Manifest 配置**
  - 应用名称、图标、主题色
  - 快捷方式（课程、进度）
  - 独立显示模式（standalone）

- **文件位置**: `public/manifest.json`

- **Vite PWA 插件配置**
  - 自动更新策略
  - 资源预缓存
  - 运行时缓存策略

- **文件位置**: `vite.config.ts`

### 3. Service Worker 离线缓存

- **缓存策略**
  - 静态资源：Cache First（缓存优先）
  - API 请求：Network First（网络优先）
  - 图片资源：Cache First + 过期策略
  - 字体资源：Cache First（长期缓存）

- **功能支持**
  - 离线访问
  - 后台同步
  - 推送通知
  - 版本更新检测

- **文件位置**: `src/service-worker.ts`

### 4. 移动端手势支持

- **支持手势**
  - 滑动（左/右/上/下）
  - 双击
  - 长按
  - 捏合缩放

- **自定义 Hooks**
  - `useSwipe` - 滑动手势
  - `useDoubleTap` - 双击手势
  - `useLongPress` - 长按手势
  - `usePinch` - 捏合手势

- **文件位置**: `src/hooks/useGestures.ts`

### 5. 移动端导航（底部导航栏）

- **功能特性**
  - 滚动自动隐藏/显示
  - 活动状态指示器
  - iOS 安全区域适配
  - 深色模式支持
  - 触摸反馈效果

- **导航项**
  - 首页
  - 课程
  - 练习
  - 进度
  - 排行榜

- **文件位置**: 
  - `src/components/MobileNav.tsx`
  - `src/components/MobileNav.css`

## 🚀 附加功能

### PWA 安装提示

- 自动检测安装能力
- 友好的安装引导
- iOS Safari 特别说明
- 支持深色模式

**文件位置**: 
- `src/components/PWAInstallPrompt.tsx`
- `src/components/PWAInstallPrompt.css`

### 离线页面

- 网络状态检测
- 友好的离线提示
- 快速重连按钮
- 离线功能说明

**文件位置**: 
- `src/components/OfflinePage.tsx`
- `src/components/OfflinePage.css`

## 📋 使用指南

### 添加到主屏幕

#### Android (Chrome)
1. 访问网站
2. 点击浏览器菜单
3. 选择"添加到主屏幕"
4. 确认添加

#### iOS (Safari)
1. 访问网站
2. 点击底部"分享"按钮
3. 选择"添加到主屏幕"
4. 确认添加

### 测试 PWA

```bash
# 开发环境
npm run dev

# 构建生产版本
npm run build

# 预览生产版本
npm run preview
```

### Lighthouse 测试

使用 Chrome DevTools 的 Lighthouse 工具测试 PWA 指标：
- PWA 得分
- 性能得分
- 可访问性得分
- 最佳实践得分

## 🎯 性能优化

### 代码分割
- React 相关库单独打包
- UI 库单独打包
- 编辑器库单独打包

### 资源优化
- 图片懒加载
- 字体子集化
- CSS 压缩
- JS 压缩

### 缓存优化
- 静态资源长期缓存
- API 响应短期缓存
- 版本更新自动检测

## 📊 移动端指标

- 首屏加载时间：< 3s
- 可交互时间：< 5s
- 触摸响应延迟：< 100ms
- 滚动帧率：≥ 60fps

## 🔧 配置说明

### vite.config.ts

```typescript
VitePWA({
  registerType: 'autoUpdate',  // 自动更新
  workbox: {
    cleanupOutdatedCaches: true,  // 清理旧缓存
    clientsClaim: true,           // 立即接管
    skipWaiting: true,            // 跳过等待
  }
})
```

### index.html

```html
<!-- 关键 meta 标签 -->
<meta name="viewport" content="width=device-width, initial-scale=1.0, viewport-fit=cover, user-scalable=no" />
<meta name="theme-color" content="#4f46e5" />
<meta name="apple-mobile-web-app-capable" content="yes" />
<link rel="manifest" href="/manifest.json" />
```

## 🐛 已知问题

1. **iOS 推送通知**: iOS Safari 不支持 Web Push API
2. **后台同步**: iOS 不支持 Background Sync API
3. **文件访问**: 移动端文件系统访问受限

## 📝 待办事项

- [ ] 添加更多页面到 PWA 快捷方式
- [ ] 实现离线数据同步队列
- [ ] 添加启动画面（splash screen）
- [ ] 优化首屏加载速度
- [ ] 添加更多手势支持（如边缘滑动）

## 📚 参考资料

- [PWA MDN 文档](https://developer.mozilla.org/zh-CN/docs/Web/Progressive_web_apps)
- [Workbox 文档](https://developers.google.com/web/tools/workbox)
- [Vite PWA 插件](https://vite-pwa-org.netlify.app/)
- [Web App Manifest](https://developer.mozilla.org/zh-CN/docs/Web/Manifest)
