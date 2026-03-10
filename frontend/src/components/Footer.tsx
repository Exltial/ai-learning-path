import { Code } from 'lucide-react'

export default function Footer() {
  return (
    <footer className="bg-secondary-800 text-white py-8 mt-auto">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {/* Brand */}
          <div>
            <div className="flex items-center space-x-2 mb-4">
              <Code className="h-6 w-6 text-primary-400" />
              <span className="text-lg font-bold">AI 学习之路</span>
            </div>
            <p className="text-secondary-300 text-sm">
              智能编程学习平台，让 AI 陪你一起成长
            </p>
          </div>

          {/* Quick Links */}
          <div>
            <h3 className="text-sm font-semibold uppercase tracking-wider mb-4">
              快速链接
            </h3>
            <ul className="space-y-2 text-secondary-300 text-sm">
              <li><a href="/courses" className="hover:text-white transition-colors">全部课程</a></li>
              <li><a href="/learning-path" className="hover:text-white transition-colors">学习路径</a></li>
              <li><a href="/profile" className="hover:text-white transition-colors">个人中心</a></li>
            </ul>
          </div>

          {/* Contact */}
          <div>
            <h3 className="text-sm font-semibold uppercase tracking-wider mb-4">
              联系我们
            </h3>
            <ul className="space-y-2 text-secondary-300 text-sm">
              <li>邮箱：support@ailearning.com</li>
              <li>微信：AI 学习助手</li>
            </ul>
          </div>
        </div>

        <div className="border-t border-secondary-700 mt-8 pt-8 text-center text-secondary-400 text-sm">
          <p>&copy; 2024 AI 学习之路。All rights reserved.</p>
        </div>
      </div>
    </footer>
  )
}
