/// <reference lib="webworker" />

import { precacheAndRoute, cleanupOutdatedCaches } from 'workbox-precaching';
import { registerRoute, NavigationRoute } from 'workbox-routing';
import { NetworkFirst, CacheFirst, StaleWhileRevalidate } from 'workbox-strategies';
import { ExpirationPlugin } from 'workbox-expiration';

// 声明 self 类型
declare const self: ServiceWorkerGlobalScope;

// 预缓存所有构建资源
precacheAndRoute(self.__WB_MANIFEST);

// 清理过时的缓存
cleanupOutdatedCaches();

// 离线缓存策略配置
const CACHE_NAMES = {
  STATIC: 'static-cache-v1',
  DYNAMIC: 'dynamic-cache-v1',
  IMAGES: 'images-cache-v1',
  API: 'api-cache-v1',
};

const CACHE_MAX_ENTRIES = 60;
const CACHE_MAX_AGE_SECONDS = 7 * 24 * 60 * 60; // 7 天

// 导航请求 - 网络优先，离线时返回缓存
const navigationHandler = new NetworkFirst({
  cacheName: CACHE_NAMES.STATIC,
  networkTimeoutSeconds: 3,
  plugins: [
    new ExpirationPlugin({
      maxEntries: CACHE_MAX_ENTRIES,
      maxAgeSeconds: CACHE_MAX_AGE_SECONDS,
    }),
  ],
});

registerRoute(
  ({ request }) => request.mode === 'navigate',
  navigationHandler
);

// 静态资源 - 缓存优先
const staticAssetHandler = new CacheFirst({
  cacheName: CACHE_NAMES.STATIC,
  plugins: [
    new ExpirationPlugin({
      maxEntries: CACHE_MAX_ENTRIES,
      maxAgeSeconds: 30 * 24 * 60 * 60, // 30 天
    }),
  ],
});

registerRoute(
  ({ request, url }) =>
    request.destination === 'script' ||
    request.destination === 'style' ||
    url.pathname.startsWith('/assets/'),
  staticAssetHandler
);

// 图片资源 - 缓存优先
const imageHandler = new CacheFirst({
  cacheName: CACHE_NAMES.IMAGES,
  plugins: [
    new ExpirationPlugin({
      maxEntries: 60,
      maxAgeSeconds: 30 * 24 * 60 * 60,
    }),
  ],
});

registerRoute(
  ({ request }) => request.destination === 'image',
  imageHandler
);

// API 请求 - 网络优先，离线时返回缓存
const apiHandler = new NetworkFirst({
  cacheName: CACHE_NAMES.API,
  networkTimeoutSeconds: 5,
  plugins: [
    new ExpirationPlugin({
      maxEntries: CACHE_MAX_ENTRIES,
      maxAgeSeconds: 24 * 60 * 60, // 1 天
    }),
  ],
});

registerRoute(
  ({ url }) => url.pathname.startsWith('/api/'),
  apiHandler
);

// 字体资源 - 缓存优先
const fontHandler = new CacheFirst({
  cacheName: CACHE_NAMES.STATIC,
  plugins: [
    new ExpirationPlugin({
      maxEntries: 30,
      maxAgeSeconds: 365 * 24 * 60 * 60, // 1 年
    }),
  ],
});

registerRoute(
  ({ request }) => request.destination === 'font',
  fontHandler
);

// 监听安装事件
self.addEventListener('install', (event) => {
  console.log('[Service Worker] Installing...');
  self.skipWaiting();
});

// 监听激活事件
self.addEventListener('activate', (event) => {
  console.log('[Service Worker] Activated');
  self.clients.claim();
});

// 监听消息事件
self.addEventListener('message', (event) => {
  if (event.data && event.data.type === 'SKIP_WAITING') {
    self.skipWaiting();
  }
});

// 后台同步
self.addEventListener('sync', (event) => {
  if (event.tag === 'sync-learning-progress') {
    event.waitUntil(syncLearningProgress());
  }
});

async function syncLearningProgress() {
  // 离线时的学习进度同步逻辑
  console.log('[Service Worker] Syncing learning progress...');
}

// 推送通知
self.addEventListener('push', (event) => {
  if (event.data) {
    const data = event.data.json();
    const options = {
      body: data.body || '您有一条新通知',
      icon: '/vite.svg',
      badge: '/vite.svg',
      vibrate: [100, 50, 100],
      data: {
        dateOfArrival: Date.now(),
        primaryKey: 1,
      },
    };
    event.waitUntil(self.registration.showNotification(data.title || 'AI 学习平台', options));
  }
});

// 通知点击处理
self.addEventListener('notificationclick', (event) => {
  event.notification.close();
  event.waitUntil(
    self.clients.matchAll({ type: 'window' }).then((clientList) => {
      for (const client of clientList) {
        if (client.url === '/' && 'focus' in client) {
          return client.focus();
        }
      }
      if (self.clients.openWindow) {
        return self.clients.openWindow('/');
      }
    })
  );
});
