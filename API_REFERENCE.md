# API 端点快速参考

## 学习进度追踪 API

### 基础 URL
```
/api/v1/progress
```

### 认证
所有端点需要 JWT 认证，通过 Header 传递：
```
Authorization: Bearer {your_jwt_token}
```

---

## 端点列表

### 1. 获取课程进度
```http
GET /api/v1/progress/courses/{course_id}
```

**路径参数:**
- `course_id` (uuid, required) - 课程 ID

**响应示例:**
```json
{
  "success": true,
  "data": {
    "course_id": "uuid",
    "enrollment_id": "uuid",
    "progress_percentage": 45.5,
    "video_progress_percentage": 38.2,
    "completed_lessons": 11,
    "total_lessons": 24,
    "lessons": [
      {
        "lesson_id": "uuid",
        "title": "第一章：Python 基础",
        "is_completed": true,
        "completed_at": "2024-03-10T10:30:00Z",
        "video_position": 600,
        "video_duration": 600
      }
    ],
    "last_accessed_at": "2024-03-10T10:30:00Z"
  }
}
```

---

### 2. 更新视频进度
```http
PUT /api/v1/progress/video
```

**请求体:**
```json
{
  "lesson_id": "uuid",
  "position": 120,
  "duration": 600
}
```

**字段说明:**
- `lesson_id` (uuid, required) - 课时 ID
- `position` (int, required) - 当前播放位置（秒）
- `duration` (int, optional) - 视频总时长（秒）

**响应示例:**
```json
{
  "success": true,
  "message": "Video progress updated successfully"
}
```

**使用建议:**
- 前端应每 15-30 秒调用一次此接口
- 视频暂停或结束时也应调用
- 当 position >= duration * 0.9 时自动标记为完成

---

### 3. 标记课时完成
```http
POST /api/v1/progress/lessons/{lesson_id}/complete
```

**路径参数:**
- `lesson_id` (uuid, required) - 课时 ID

**响应示例:**
```json
{
  "success": true,
  "message": "Lesson marked as completed"
}
```

---

### 4. 获取学习热力图
```http
GET /api/v1/progress/heatmap?months=6
```

**查询参数:**
- `months` (int, optional, default: 6) - 月份数 (1-24)

**响应示例:**
```json
{
  "success": true,
  "data": {
    "heatmap": [
      {
        "date": "2024-03-10",
        "count": 45,
        "level": 2
      }
    ],
    "months": 6
  }
}
```

**热力图等级:**
- Level 0: 无活动
- Level 1: 轻度 (1-15 分钟)
- Level 2: 中度 (15-30 分钟)
- Level 3: 重度 (30-60 分钟)
- Level 4: 高强度 (60+ 分钟)

---

### 5. 获取每日统计
```http
GET /api/v1/progress/daily-stats?days=30
```

**查询参数:**
- `days` (int, optional, default: 30) - 天数 (1-365)

**响应示例:**
```json
{
  "success": true,
  "data": {
    "daily_stats": [
      {
        "date": "2024-03-10",
        "total_seconds": 2700,
        "lessons_completed": 2,
        "courses_accessed": 1
      }
    ],
    "days": 30
  }
}
```

---

### 6. 获取周报
```http
GET /api/v1/progress/reports/weekly?offset=0
```

**查询参数:**
- `offset` (int, optional, default: 0) - 周偏移 (0=本周，-1=上周，以此类推)

**响应示例:**
```json
{
  "success": true,
  "data": {
    "week_start": "2024-03-04",
    "week_end": "2024-03-10",
    "total_hours": 12.5,
    "lessons_completed": 8,
    "courses_progress": [
      {
        "course_id": "uuid",
        "course_title": "Python 编程基础",
        "progress_percent": 45.0,
        "lessons_completed": 8,
        "total_lessons": 24,
        "time_spent_minutes": 450
      }
    ],
    "daily_stats": [...],
    "avg_daily_minutes": 107.1
  }
}
```

---

### 7. 获取月报
```http
GET /api/v1/progress/reports/monthly?offset=0
```

**查询参数:**
- `offset` (int, optional, default: 0) - 月偏移 (0=本月，-1=上月，以此类推)

**响应示例:**
```json
{
  "success": true,
  "data": {
    "month": "March",
    "year": 2024,
    "total_hours": 48.5,
    "lessons_completed": 32,
    "courses_completed": 1,
    "courses_progress": [...],
    "daily_stats": [...],
    "avg_daily_minutes": 95.3,
    "best_day": "2024-03-15",
    "best_day_minutes": 185
  }
}
```

---

### 8. 导出报告
```http
POST /api/v1/progress/reports/export
```

**请求体:**
```json
{
  "report_type": "monthly",
  "offset": 0,
  "format": "csv"
}
```

**字段说明:**
- `report_type` (string, required) - 报告类型 ("weekly" 或 "monthly")
- `offset` (int, optional, default: 0) - 偏移量
- `format` (string, optional, default: "csv") - 导出格式 ("csv" 或 "json")

**响应:**
- Content-Type: `text/csv` 或 `application/json`
- Content-Disposition: `attachment; filename=learning_report_{type}_{date}.{format}`

---

### 9. 获取学习统计
```http
GET /api/v1/progress/stats
```

**响应示例:**
```json
{
  "success": true,
  "data": {
    "total_learning_seconds": 172800,
    "total_learning_hours": 48.0,
    "total_lessons": 50,
    "completed_lessons": 32,
    "completion_rate": 64.0
  }
}
```

---

## 错误响应

### 400 Bad Request
```json
{
  "success": false,
  "error": {
    "code": "BAD_REQUEST",
    "message": "Invalid course ID"
  }
}
```

### 401 Unauthorized
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "User ID not found"
  }
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "Error details"
  }
}
```

---

## 前端集成示例

### React - 更新视频进度
```tsx
// 在视频播放器组件中
const handleTimeUpdate = () => {
  const position = Math.floor(videoRef.current.currentTime)
  
  // 每 30 秒保存一次
  if (position % 30 === 0 && position !== lastSavedPosition) {
    api.put('/api/v1/progress/video', {
      lesson_id: lessonId,
      position: position,
      duration: videoDuration
    })
    lastSavedPosition = position
  }
}

// 视频结束时标记完成
const handleVideoEnded = () => {
  api.post(`/api/v1/progress/lessons/${lessonId}/complete`)
}
```

### React - 获取热力图数据
```tsx
const [heatmapData, setHeatmapData] = useState([])

useEffect(() => {
  const loadHeatmap = async () => {
    const response = await api.get('/api/v1/progress/heatmap?months=6')
    if (response.success) {
      setHeatmapData(response.data.heatmap)
    }
  }
  loadHeatmap()
}, [])
```

### React - 导出报告
```tsx
const handleExport = async () => {
  const response = await api.post('/api/v1/progress/reports/export', {
    report_type: 'monthly',
    offset: 0,
    format: 'csv'
  }, {
    responseType: 'blob'
  })
  
  const blob = new Blob([response.data], { type: 'text/csv' })
  const url = window.URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = 'learning_report.csv'
  link.click()
  window.URL.revokeObjectURL(url)
}
```

---

## 最佳实践

1. **视频进度保存频率**: 建议每 15-30 秒调用一次更新接口
2. **错误处理**: 所有 API 调用都应该有 try-catch 错误处理
3. **加载状态**: 显示加载指示器提升用户体验
4. **数据缓存**: 热力图和报告数据可以适当缓存（5-10 分钟）
5. **离线支持**: 考虑在离线时本地保存进度，网络恢复后同步

---

**文档版本**: 1.0  
**更新时间**: 2026-03-10
