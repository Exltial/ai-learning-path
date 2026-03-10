# 性能优化实现指南

## 概述

本指南介绍 AI 学习平台项目的性能优化功能实现，包括：
- ✅ Redis 缓存策略
- ✅ 数据库查询优化（索引）
- ✅ 图片懒加载
- ✅ 代码分割（Code Splitting）
- ✅ 性能监控

## 性能目标

- **API 响应时间**: < 200ms
- **首屏加载时间**: < 2 秒
- **P95 响应时间**: < 300ms
- **LCP (最大内容绘制)**: < 2.5 秒

---

## 1. Redis 缓存策略

### 文件位置
`backend/internal/services/cache_service.go`

### 功能特性

- 课程数据缓存（TTL: 15 分钟）
- 用户数据缓存（TTL: 10 分钟）
- 列表查询缓存（TTL: 5 分钟）
- 自动缓存预热
- 缓存失效管理

### 使用方法

#### 初始化缓存服务

```go
import (
    "github.com/redis/go-redis/v9"
    "ai-learning-platform/internal/services"
)

// 创建 Redis 客户端
redisClient := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "", // 如果有密码
    DB:       0,
})

// 创建缓存服务
cacheService := services.NewCacheService(redisClient, nil) // nil 使用默认配置
```

#### 缓存课程数据

```go
// 获取课程（带缓存）
func (h *CourseHandler) GetCourse(c *gin.Context) {
    courseID := c.Param("id")
    
    var course models.Course
    fromCache, err := h.cacheService.GetOrSet(
        c.Request.Context(),
        h.cacheService.CacheKey("course", courseID),
        &course,
        func() (interface{}, error) {
            // 从数据库获取
            return h.courseService.GetCourse(c.Request.Context(), uuid.MustParse(courseID))
        },
        15*time.Minute,
    )
    
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{
        "success": true,
        "data": course,
        "from_cache": fromCache,
    })
}
```

#### 缓存用户数据

```go
// 获取用户（带缓存）
func (s *UserService) GetUserWithCache(ctx context.Context, userID uuid.UUID) (*models.User, error) {
    var user models.User
    _, err := s.cacheService.GetOrSet(
        ctx,
        s.cacheService.CacheKey("user", userID.String()),
        &user,
        func() (interface{}, error) {
            return s.userRepo.GetByID(ctx, userID)
        },
        10*time.Minute,
    )
    return &user, err
}
```

#### 缓存失效

```go
// 更新课程后使缓存失效
func (s *CourseService) UpdateCourse(ctx context.Context, course *models.Course) error {
    // 更新数据库
    if err := s.courseRepo.Update(ctx, course); err != nil {
        return err
    }
    
    // 使缓存失效
    s.cacheService.InvalidateCourse(ctx, course.ID.String())
    s.cacheService.InvalidateCourseList(ctx)
    
    return nil
}
```

### 缓存键命名规范

```
ai-learning:course:{id}          - 课程数据
ai-learning:user:{id}            - 用户数据
ai-learning:courses:list:{params} - 课程列表
```

---

## 2. 数据库查询优化

### 文件位置
`backend/migrations/008_add_indexes.sql`

### 索引策略

#### 用户表索引
```sql
-- 认证查询优化
CREATE INDEX idx_users_email ON users(email);

-- 常用查询组合
CREATE INDEX idx_users_status_created ON users(status, created_at);
```

#### 课程表索引
```sql
-- 筛选查询优化
CREATE INDEX idx_courses_category ON courses(category);
CREATE INDEX idx_courses_difficulty ON courses(difficulty_level);
CREATE INDEX idx_courses_is_published ON courses(is_published);

-- 组合索引（最重要）
CREATE INDEX idx_courses_published_category_difficulty 
ON courses(is_published, category, difficulty_level);
```

#### 选课记录索引
```sql
-- 防止重复 + 快速查询
CREATE UNIQUE INDEX idx_enrollments_user_course_unique 
ON enrollments(user_id, course_id);

-- 用户选课列表
CREATE INDEX idx_enrollments_user_status ON enrollments(user_id, status);
```

### 查询优化建议

#### 使用 EXPLAIN 分析查询
```sql
EXPLAIN ANALYZE 
SELECT * FROM courses 
WHERE is_published = true 
  AND category = 'programming'
  AND difficulty_level = 'beginner';
```

#### 避免 N+1 查询
```go
// ❌ 不好的做法
for _, course := range courses {
    lessons, _ := lessonRepo.GetByCourseID(ctx, course.ID)
    course.Lessons = lessons
}

// ✅ 好的做法 - 批量加载
courseIDs := make([]uuid.UUID, len(courses))
for i, c := range courses {
    courseIDs[i] = c.ID
}
lessonsMap := lessonRepo.GetByCourseIDs(ctx, courseIDs)
```

#### 使用连接池
```go
// 在 main.go 中配置
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

---

## 3. 图片懒加载

### 文件位置
`frontend/src/hooks/useLazyLoad.ts`

### 功能特性

- Intersection Observer API
- 视口外预加载（loadOffset）
- 占位图支持
- 错误处理
- 平滑过渡动画

### 使用方法

#### 懒加载图片组件

```tsx
import { LazyImage } from '@/hooks/useLazyLoad'

function CourseCard({ course }) {
  return (
    <div className="course-card">
      <LazyImage
        src={course.thumbnail}
        alt={course.title}
        placeholder="/images/placeholder.jpg"
        threshold={0.1}
        rootMargin="100px"
        className="course-thumbnail"
      />
      <h3>{course.title}</h3>
    </div>
  )
}
```

#### 自定义懒加载 Hook

```tsx
import { useLazyLoad } from '@/hooks/useLazyLoad'

function HeavyComponent() {
  const { ref, isVisible, isLoaded } = useLazyLoad({
    threshold: 0,
    rootMargin: '50px',
  })

  return (
    <div ref={ref}>
      {isVisible && (
        <div style={{ opacity: isLoaded ? 1 : 0, transition: 'opacity 0.3s' }}>
          {/* 重型组件内容 */}
        </div>
      )}
    </div>
  )
}
```

#### 懒加载 React 组件

```tsx
import { lazyComponent } from '@/hooks/useLazyLoad'

// 重型组件
function AnalyticsDashboard() {
  // 复杂的图表渲染逻辑
  return <div>...</div>
}

// 创建懒加载版本
const LazyAnalytics = lazyComponent(AnalyticsDashboard, {
  threshold: 0.1,
  fallback: <div>Loading dashboard...</div>,
})

// 使用
function Page() {
  return (
    <div>
      <LazyAnalytics />
    </div>
  )
}
```

---

## 4. 代码分割（Code Splitting）

### 实现方式

#### 路由级代码分割

```tsx
// App.tsx
import { Suspense, lazy } from 'react'
import { BrowserRouter, Routes, Route } from 'react-router-dom'

// 懒加载页面
const HomePage = lazy(() => import('@/pages/HomePage'))
const CoursesPage = lazy(() => import('@/pages/CoursesPage'))
const CourseDetailPage = lazy(() => import('@/pages/CourseDetailPage'))
const ProfilePage = lazy(() => import('@/pages/ProfilePage'))

function App() {
  return (
    <BrowserRouter>
      <Suspense fallback={<div className="loading-spinner">Loading...</div>}>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/courses" element={<CoursesPage />} />
          <Route path="/courses/:id" element={<CourseDetailPage />} />
          <Route path="/profile" element={<ProfilePage />} />
        </Routes>
      </Suspense>
    </BrowserRouter>
  )
}
```

#### 组件级代码分割

```tsx
// 懒加载重型组件
const CodeEditor = lazy(() => import('@/components/CodeEditor'))
const LearningHeatmap = lazy(() => import('@/components/LearningHeatmap'))

function CourseContent() {
  return (
    <Suspense fallback={<div>Loading editor...</div>}>
      <CodeEditor />
      <LearningHeatmap />
    </Suspense>
  )
}
```

#### 动态导入

```tsx
// 按需加载
async function loadExercise(type: string) {
  switch (type) {
    case 'coding':
      return import('@/exercises/CodingExercise')
    case 'quiz':
      return import('@/exercises/QuizExercise')
    case 'project':
      return import('@/exercises/ProjectExercise')
  }
}
```

### Vite 配置优化

```ts
// vite.config.ts
export default defineConfig({
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          'react-vendor': ['react', 'react-dom', 'react-router-dom'],
          'ui-vendor': ['lucide-react'],
          'editor-vendor': ['@monaco-editor/react'],
        },
      },
    },
  },
})
```

---

## 5. 性能监控

### 后端监控

#### 文件位置
`backend/internal/middleware/performance_middleware.go`

#### 使用方法

```go
// main.go
import (
    "ai-learning-platform/internal/middleware"
)

func main() {
    r := gin.Default()
    
    // 获取全局性能指标
    metrics := middleware.GetGlobalMetrics()
    
    // 应用性能中间件
    r.Use(middleware.PerformanceMiddleware(metrics))
    
    // 性能监控端点
    r.GET("/api/performance/metrics", middleware.PerformanceMonitorHandler(metrics))
    
    // ... 其他路由
}
```

#### 监控指标

- 总请求数
- 慢请求数（> 200ms）
- 平均响应时间
- P95/P99 响应时间
- 最大/最小响应时间

### 前端监控

#### 文件位置
`frontend/src/components/PerformanceMonitor.tsx`

#### 使用方法

```tsx
// App.tsx
import { PerformanceMonitor } from '@/components/PerformanceMonitor'

function App() {
  return (
    <>
      <PerformanceMonitor 
        apiMetricsEndpoint="/api/performance/metrics"
        refreshInterval={5000}
        position="bottom-right"
      />
      {/* 其他应用内容 */}
    </>
  )
}
```

#### 监控指标

**API 性能:**
- 平均响应时间
- P95/P99 响应时间
- 慢请求百分比

**页面加载:**
- FCP (首次内容绘制)
- LCP (最大内容绘制)
- TTI (可交互时间)
- CLS (累积布局偏移)

#### 快捷键

- `Shift + P`: 切换监控面板显示

### 性能预算

```go
// 设置性能预算中间件
r.Use(middleware.PerformanceBudgetMiddleware(metrics, 200*time.Millisecond))
```

---

## 集成示例

### 完整的课程服务（带缓存）

```go
package services

type CourseService struct {
    courseRepo     *repository.CourseRepository
    cacheService   *services.CacheService
}

func (s *CourseService) ListCourses(ctx context.Context, category, difficulty string, page, limit int) ([]*models.Course, int, error) {
    // 尝试从缓存获取
    cached, err := s.cacheService.GetCourseList(ctx, category, difficulty, page, limit)
    if err == nil && cached != nil {
        return cached.(*[]*models.Course), 0, nil
    }
    
    // 从数据库获取
    courses, total, err := s.courseRepo.List(ctx, category, difficulty, page, limit)
    if err != nil {
        return nil, 0, err
    }
    
    // 写入缓存
    s.cacheService.SetCourseList(ctx, category, difficulty, page, limit, courses)
    
    return courses, total, nil
}
```

### 完整的课程卡片组件（带懒加载）

```tsx
import { LazyImage } from '@/hooks/useLazyLoad'
import { Link } from 'react-router-dom'

interface CourseCardProps {
  course: Course
}

export function CourseCard({ course }: CourseCardProps) {
  return (
    <div className="course-card">
      <LazyImage
        src={course.thumbnail}
        alt={course.title}
        placeholder="/images/course-placeholder.jpg"
        threshold={0.1}
        rootMargin="100px"
      />
      
      <div className="course-info">
        <h3>{course.title}</h3>
        <p className="course-description">{course.description}</p>
        
        <div className="course-meta">
          <span className="difficulty">{course.difficulty_level}</span>
          <span className="students">{course.enrollment_count} 学生</span>
          <span className="rating">⭐ {course.rating}</span>
        </div>
        
        <Link to={`/courses/${course.id}`} className="course-link">
          查看详情
        </Link>
      </div>
    </div>
  )
}
```

---

## 性能测试

### 后端性能测试

```bash
# 使用 ab 进行压力测试
ab -n 1000 -c 10 http://localhost:8080/api/courses

# 使用 wrk
wrk -t12 -c400 -d30s http://localhost:8080/api/courses
```

### 前端性能测试

```bash
# 使用 Lighthouse CI
npx @lhci/cli@0.11.x autorun

# 使用 WebPageTest
# 访问 https://www.webpagetest.org/
```

### 性能检查清单

- [ ] API 响应时间 < 200ms
- [ ] 首屏加载 < 2 秒
- [ ] LCP < 2.5 秒
- [ ] FID < 100ms
- [ ] CLS < 0.1
- [ ] 图片实现懒加载
- [ ] 路由实现代码分割
- [ ] 数据库查询使用索引
- [ ] 热点数据使用缓存

---

## 故障排查

### 缓存问题

```go
// 检查 Redis 连接
if !cacheService.HealthCheck(ctx) {
    log.Error("Redis connection failed")
}

// 查看缓存统计
stats, _ := cacheService.GetCacheStats(ctx)
log.Info("Cache stats: %+v", stats)
```

### 性能下降排查

1. 检查性能监控面板
2. 查看慢查询日志
3. 分析数据库 EXPLAIN 输出
4. 检查 Redis 命中率
5. 查看前端 Network 面板

---

## 最佳实践

1. **缓存策略**: 对读多写少的数据使用缓存
2. **索引设计**: 根据实际查询模式设计索引
3. **懒加载**: 对首屏外的内容使用懒加载
4. **代码分割**: 按路由和功能模块分割代码
5. **监控告警**: 设置性能预算和告警阈值

---

## 相关文件

- `backend/internal/services/cache_service.go` - Redis 缓存服务
- `backend/internal/middleware/performance_middleware.go` - 性能监控中间件
- `backend/migrations/008_add_indexes.sql` - 数据库索引迁移
- `frontend/src/hooks/useLazyLoad.ts` - 懒加载 Hook
- `frontend/src/components/PerformanceMonitor.tsx` - 性能监控组件

---

*最后更新：2026-03-10*
