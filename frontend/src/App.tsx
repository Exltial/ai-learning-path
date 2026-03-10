import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { AuthProvider } from './contexts/AuthContext'
import Layout from './components/Layout'
import OfflinePage from './components/OfflinePage'
import HomePage from './pages/HomePage'
import CoursesPage from './pages/CoursesPage'
import CourseDetailPage from './pages/CourseDetailPage'
import LearningPathPage from './pages/LearningPathPage'
import ProfilePage from './pages/ProfilePage'
import DiscussionPage from './pages/DiscussionPage'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'

// 简单占位页面组件
const PracticePage = () => (
  <div className="container mx-auto px-4 py-8 page-transition">
    <h1 className="text-2xl font-bold mb-4">代码练习</h1>
    <p className="text-gray-600 dark:text-gray-400">
      编程练习功能开发中...
    </p>
  </div>
)

const ProgressPage = () => (
  <div className="container mx-auto px-4 py-8 page-transition">
    <h1 className="text-2xl font-bold mb-4">学习进度</h1>
    <p className="text-gray-600 dark:text-gray-400">
      学习进度追踪功能开发中...
    </p>
  </div>
)

const LeaderboardPage = () => (
  <div className="container mx-auto px-4 py-8 page-transition">
    <h1 className="text-2xl font-bold mb-4">排行榜</h1>
    <p className="text-gray-600 dark:text-gray-400">
      排行榜功能开发中...
    </p>
  </div>
)

function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <OfflinePage />
        <Routes>
          {/* Auth Routes (no layout) */}
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />

          {/* Main Routes (with layout) */}
          <Route path="/" element={<Layout />}>
            <Route index element={<HomePage />} />
            <Route path="courses" element={<CoursesPage />} />
            <Route path="courses/:id" element={<CourseDetailPage />} />
            <Route path="courses/:courseId/discussions" element={<DiscussionPage />} />
            <Route path="learning-path" element={<LearningPathPage />} />
            <Route path="practice" element={<PracticePage />} />
            <Route path="progress" element={<ProgressPage />} />
            <Route path="leaderboard" element={<LeaderboardPage />} />
            <Route path="profile" element={<ProfilePage />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  )
}

export default App
