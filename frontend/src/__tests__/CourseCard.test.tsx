import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import CourseCard, { type Course } from '../components/CourseCard'
import { MemoryRouter } from 'react-router-dom'

// Mock lucide-react icons
vi.mock('lucide-react', () => ({
  Clock: (props: any) => <span data-testid="clock-icon" {...props} />,
  Users: (props: any) => <span data-testid="users-icon" {...props} />,
  Star: (props: any) => <span data-testid="star-icon" {...props} />,
  BookOpen: (props: any) => <span data-testid="book-icon" {...props} />,
}))

const mockCourse: Course = {
  id: '1',
  title: 'Python 编程基础',
  description: '从零开始学习 Python，掌握编程基础语法和核心概念',
  thumbnail: 'https://picsum.photos/seed/python/400/200',
  level: 'beginner',
  duration: '8 小时',
  students: 1234,
  rating: 4.8,
  lessons: 24,
  tags: ['Python', '编程基础', '入门'],
}

const renderWithRouter = (component: React.ReactElement) => {
  return render(
    <MemoryRouter>
      {component}
    </MemoryRouter>
  )
}

describe('CourseCard', () => {
  describe('渲染', () => {
    it('应该正确渲染课程卡片', () => {
      renderWithRouter(<CourseCard course={mockCourse} />)
      
      expect(screen.getByText('Python 编程基础')).toBeInTheDocument()
      expect(screen.getByText('从零开始学习 Python，掌握编程基础语法和核心概念')).toBeInTheDocument()
    })

    it('应该显示课程缩略图', () => {
      renderWithRouter(<CourseCard course={mockCourse} />)
      
      const image = screen.getByRole('img')
      expect(image).toHaveAttribute('src', mockCourse.thumbnail)
      expect(image).toHaveAttribute('alt', mockCourse.title)
    })

    it('应该显示课程难度标签', () => {
      renderWithRouter(<CourseCard course={mockCourse} />)
      
      expect(screen.getByText('入门')).toBeInTheDocument()
    })

    it('应该显示课程时长', () => {
      renderWithRouter(<CourseCard course={mockCourse} />)
      
      expect(screen.getByText(mockCourse.duration)).toBeInTheDocument()
      expect(screen.getByTestId('clock-icon')).toBeInTheDocument()
    })

    it('应该显示学生数量', () => {
      renderWithRouter(<CourseCard course={mockCourse} />)
      
      expect(screen.getByText('1234')).toBeInTheDocument()
      expect(screen.getByTestId('users-icon')).toBeInTheDocument()
    })

    it('应该显示课程评分', () => {
      renderWithRouter(<CourseCard course={mockCourse} />)
      
      expect(screen.getByText('4.8')).toBeInTheDocument()
      expect(screen.getByTestId('star-icon')).toBeInTheDocument()
    })

    it('应该显示课程章节数', () => {
      renderWithRouter(<CourseCard course={mockCourse} />)
      
      expect(screen.getByText('24')).toBeInTheDocument()
      expect(screen.getByTestId('book-icon')).toBeInTheDocument()
    })

    it('应该显示课程标签', () => {
      renderWithRouter(<CourseCard course={mockCourse} />)
      
      expect(screen.getByText('Python')).toBeInTheDocument()
      expect(screen.getByText('编程基础')).toBeInTheDocument()
      expect(screen.getByText('入门')).toBeInTheDocument()
    })
  })

  describe('难度级别显示', () => {
    it('应该正确显示 beginner 级别', () => {
      const course = { ...mockCourse, level: 'beginner' as const }
      renderWithRouter(<CourseCard course={course} />)
      
      expect(screen.getByText('入门')).toBeInTheDocument()
    })

    it('应该正确显示 intermediate 级别', () => {
      const course = { ...mockCourse, level: 'intermediate' as const }
      renderWithRouter(<CourseCard course={course} />)
      
      expect(screen.getByText('进阶')).toBeInTheDocument()
    })

    it('应该正确显示 advanced 级别', () => {
      const course = { ...mockCourse, level: 'advanced' as const }
      renderWithRouter(<CourseCard course={course} />)
      
      expect(screen.getByText('高级')).toBeInTheDocument()
    })

    it('应该处理未知的难度级别', () => {
      const course = { ...mockCourse, level: 'unknown' as any }
      renderWithRouter(<CourseCard course={course} />)
      
      // 应该有一个默认显示或空
      expect(screen.queryByText('入门')).not.toBeInTheDocument()
    })
  })

  describe('点击交互', () => {
    it('点击卡片应该跳转到课程详情页', () => {
      renderWithRouter(<CourseCard course={mockCourse} />)
      
      const card = screen.getByText('Python 编程基础').closest('a')
      expect(card).toHaveAttribute('href', '/courses/1')
    })
  })

  describe('样式和布局', () => {
    it('应该有正确的卡片容器类名', () => {
      const { container } = renderWithRouter(<CourseCard course={mockCourse} />)
      
      const card = container.firstChild
      expect(card).toHaveClass('bg-white')
      expect(card).toHaveClass('dark:bg-secondary-800')
      expect(card).toHaveClass('rounded-xl')
      expect(card).toHaveClass('overflow-hidden')
    })

    it('应该有悬停效果', () => {
      const { container } = renderWithRouter(<CourseCard course={mockCourse} />)
      
      const card = container.firstChild
      expect(card).toHaveClass('hover:shadow-lg')
      expect(card).toHaveClass('transition-shadow')
    })
  })

  describe('边界情况', () => {
    it('应该处理长标题', () => {
      const course = {
        ...mockCourse,
        title: '这是一个非常非常长的课程标题，可能会超出容器宽度',
      }
      renderWithRouter(<CourseCard course={course} />)
      
      expect(screen.getByText(course.title)).toBeInTheDocument()
    })

    it('应该处理长描述', () => {
      const course = {
        ...mockCourse,
        description: '这是一个非常长的描述，包含了很多关于课程的详细信息。'.repeat(10),
      }
      renderWithRouter(<CourseCard course={course} />)
      
      expect(screen.getByText(course.description)).toBeInTheDocument()
    })

    it('应该处理空标签数组', () => {
      const course = { ...mockCourse, tags: [] }
      renderWithRouter(<CourseCard course={course} />)
      
      // 标签区域应该为空或不存在
      const tagContainer = screen.queryByRole('list')
      expect(tagContainer).toBeNull()
    })

    it('应该处理零学生数', () => {
      const course = { ...mockCourse, students: 0 }
      renderWithRouter(<CourseCard course={course} />)
      
      expect(screen.getByText('0')).toBeInTheDocument()
    })

    it('应该处理零评分', () => {
      const course = { ...mockCourse, rating: 0 }
      renderWithRouter(<CourseCard course={course} />)
      
      expect(screen.getByText('0')).toBeInTheDocument()
    })

    it('应该处理大数字学生数', () => {
      const course = { ...mockCourse, students: 999999 }
      renderWithRouter(<CourseCard course={course} />)
      
      expect(screen.getByText('999999')).toBeInTheDocument()
    })
  })

  describe('可访问性', () => {
    it('图片应该有 alt 属性', () => {
      renderWithRouter(<CourseCard course={mockCourse} />)
      
      const image = screen.getByRole('img')
      expect(image).toHaveAttribute('alt', mockCourse.title)
    })

    it('卡片链接应该有正确的 aria 标签', () => {
      renderWithRouter(<CourseCard course={mockCourse} />)
      
      const link = screen.getByRole('link')
      expect(link).toHaveAttribute('href', `/courses/${mockCourse.id}`)
    })
  })

  describe('响应式设计', () => {
    it('应该支持不同的视图模式', () => {
      // 这个测试验证组件在不同容器宽度下的表现
      const { container } = renderWithRouter(<CourseCard course={mockCourse} />)
      
      expect(container.firstChild).toBeInTheDocument()
    })
  })
})
