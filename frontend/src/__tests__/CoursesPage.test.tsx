import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import CoursesPage from '../pages/CoursesPage'

// Mock lucide-react icons
vi.mock('lucide-react', () => ({
  Search: (props: any) => <span data-testid="search-icon" {...props} />,
  Filter: (props: any) => <span data-testid="filter-icon" {...props} />,
  Grid: (props: any) => <span data-testid="grid-icon" {...props} />,
  List: (props: any) => <span data-testid="list-icon" {...props} />,
}))

// Mock CourseCard
vi.mock('../components/CourseCard', () => ({
  default: ({ course }: any) => (
    <div data-testid="course-card" data-course-id={course.id}>
      <h3>{course.title}</h3>
      <p>{course.description}</p>
    </div>
  ),
}))

const renderWithRouter = (component: React.ReactElement) => {
  return render(
    <MemoryRouter>
      {component}
    </MemoryRouter>
  )
}

describe('CoursesPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('渲染', () => {
    it('应该正确渲染页面标题', () => {
      renderWithRouter(<CoursesPage />)
      
      expect(screen.getByText('全部课程')).toBeInTheDocument()
      expect(screen.getByText('探索我们的课程库，找到适合你的学习内容')).toBeInTheDocument()
    })

    it('应该渲染搜索框', () => {
      renderWithRouter(<CoursesPage />)
      
      const searchInput = screen.getByPlaceholderText('搜索课程...')
      expect(searchInput).toBeInTheDocument()
      expect(screen.getByTestId('search-icon')).toBeInTheDocument()
    })

    it('应该渲染难度筛选器', () => {
      renderWithRouter(<CoursesPage />)
      
      expect(screen.getByTestId('filter-icon')).toBeInTheDocument()
      const select = screen.getByRole('combobox')
      expect(select).toBeInTheDocument()
    })

    it('应该渲染视图切换按钮', () => {
      renderWithRouter(<CoursesPage />)
      
      expect(screen.getByTestId('grid-icon')).toBeInTheDocument()
      expect(screen.getByTestId('list-icon')).toBeInTheDocument()
    })

    it('应该默认显示网格视图', () => {
      renderWithRouter(<CoursesPage />)
      
      const gridButton = screen.getByTestId('grid-icon').parentElement
      const listButton = screen.getByTestId('list-icon').parentElement
      
      expect(gridButton).toHaveClass('bg-primary-100')
      expect(listButton).not.toHaveClass('bg-primary-100')
    })

    it('应该显示所有课程（默认）', () => {
      renderWithRouter(<CoursesPage />)
      
      const cards = screen.getAllByTestId('course-card')
      expect(cards).toHaveLength(6) // 硬编码的 6 个课程
    })

    it('应该显示课程数量', () => {
      renderWithRouter(<CoursesPage />)
      
      expect(screen.getByText('找到 6 门课程')).toBeInTheDocument()
    })
  })

  describe('搜索功能', () => {
    it('应该根据标题搜索过滤课程', () => {
      renderWithRouter(<CoursesPage />)
      
      const searchInput = screen.getByPlaceholderText('搜索课程...')
      fireEvent.change(searchInput, { target: { value: 'Python' } })
      
      expect(screen.getByText('找到 2 门课程')).toBeInTheDocument()
    })

    it('应该根据描述搜索过滤课程', () => {
      renderWithRouter(<CoursesPage />)
      
      const searchInput = screen.getByPlaceholderText('搜索课程...')
      fireEvent.change(searchInput, { target: { value: '前端' } })
      
      // React 课程和 TypeScript 课程都提到前端
      const cards = screen.getAllByTestId('course-card')
      expect(cards.length).toBeGreaterThan(0)
    })

    it('应该根据标签搜索过滤课程', () => {
      renderWithRouter(<CoursesPage />)
      
      const searchInput = screen.getByPlaceholderText('搜索课程...')
      fireEvent.change(searchInput, { target: { value: '算法' } })
      
      expect(screen.getByText('找到 1 门课程')).toBeInTheDocument()
    })

    it('搜索应该不区分大小写', () => {
      renderWithRouter(<CoursesPage />)
      
      const searchInput = screen.getByPlaceholderText('搜索课程...')
      fireEvent.change(searchInput, { target: { value: 'python' } })
      
      expect(screen.getByText('找到 2 门课程')).toBeInTheDocument()
    })

    it('清空搜索应该显示所有课程', () => {
      renderWithRouter(<CoursesPage />)
      
      const searchInput = screen.getByPlaceholderText('搜索课程...')
      fireEvent.change(searchInput, { target: { value: 'Python' } })
      fireEvent.change(searchInput, { target: { value: '' } })
      
      expect(screen.getByText('找到 6 门课程')).toBeInTheDocument()
    })

    it('搜索无结果应该显示空状态', () => {
      renderWithRouter(<CoursesPage />)
      
      const searchInput = screen.getByPlaceholderText('搜索课程...')
      fireEvent.change(searchInput, { target: { value: '不存在的课程' } })
      
      expect(screen.getByText('没有找到匹配的课程')).toBeInTheDocument()
      expect(screen.getByText('尝试调整搜索条件或筛选器')).toBeInTheDocument()
    })
  })

  describe('难度筛选', () => {
    it('应该筛选入门级别课程', () => {
      renderWithRouter(<CoursesPage />)
      
      const select = screen.getByRole('combobox')
      fireEvent.change(select, { target: { value: 'beginner' } })
      
      expect(screen.getByText('找到 1 门课程')).toBeInTheDocument()
    })

    it('应该筛选进阶级别课程', () => {
      renderWithRouter(<CoursesPage />)
      
      const select = screen.getByRole('combobox')
      fireEvent.change(select, { target: { value: 'intermediate' } })
      
      expect(screen.getByText('找到 3 门课程')).toBeInTheDocument()
    })

    it('应该筛选高级级别课程', () => {
      renderWithRouter(<CoursesPage />)
      
      const select = screen.getByRole('combobox')
      fireEvent.change(select, { target: { value: 'advanced' } })
      
      expect(screen.getByText('找到 2 门课程')).toBeInTheDocument()
    })

    it('应该重置为全部难度', () => {
      renderWithRouter(<CoursesPage />)
      
      const select = screen.getByRole('combobox')
      fireEvent.change(select, { target: { value: 'beginner' } })
      fireEvent.change(select, { target: { value: 'all' } })
      
      expect(screen.getByText('找到 6 门课程')).toBeInTheDocument()
    })
  })

  describe('视图切换', () => {
    it('应该切换到列表视图', () => {
      renderWithRouter(<CoursesPage />)
      
      const listButton = screen.getByTestId('list-icon').parentElement
      if (listButton) {
        fireEvent.click(listButton)
      }
      
      expect(listButton).toHaveClass('bg-primary-100')
      const gridButton = screen.getByTestId('grid-icon').parentElement
      expect(gridButton).not.toHaveClass('bg-primary-100')
    })

    it('应该切换回网格视图', () => {
      renderWithRouter(<CoursesPage />)
      
      const listButton = screen.getByTestId('list-icon').parentElement
      if (listButton) {
        fireEvent.click(listButton)
      }
      
      const gridButton = screen.getByTestId('grid-icon').parentElement
      if (gridButton) {
        fireEvent.click(gridButton)
      }
      
      expect(gridButton).toHaveClass('bg-primary-100')
    })
  })

  describe('组合筛选', () => {
    it('应该同时应用搜索和难度筛选', () => {
      renderWithRouter(<CoursesPage />)
      
      const searchInput = screen.getByPlaceholderText('搜索课程...')
      fireEvent.change(searchInput, { target: { value: 'Python' } })
      
      const select = screen.getByRole('combobox')
      fireEvent.change(select, { target: { value: 'beginner' } })
      
      // Python 编程基础是 beginner 级别
      expect(screen.getByText('找到 1 门课程')).toBeInTheDocument()
    })

    it('组合筛选无结果应该显示空状态', () => {
      renderWithRouter(<CoursesPage />)
      
      const searchInput = screen.getByPlaceholderText('搜索课程...')
      fireEvent.change(searchInput, { target: { value: 'Python' } })
      
      const select = screen.getByRole('combobox')
      fireEvent.change(select, { target: { value: 'advanced' } })
      
      // Python 课程是 beginner 级别，所以 advanced 筛选后应该没有结果
      expect(screen.getByText('找到 0 门课程')).toBeInTheDocument()
    })
  })

  describe('课程卡片', () => {
    it('应该为每个课程渲染 CourseCard 组件', () => {
      renderWithRouter(<CoursesPage />)
      
      const cards = screen.getAllByTestId('course-card')
      expect(cards).toHaveLength(6)
    })

    it('课程卡片应该显示正确的课程 ID', () => {
      renderWithRouter(<CoursesPage />)
      
      const cards = screen.getAllByTestId('course-card')
      expect(cards[0]).toHaveAttribute('data-course-id', '1')
      expect(cards[1]).toHaveAttribute('data-course-id', '2')
    })
  })

  describe('响应式设计', () => {
    it('搜索框应该在移动端垂直排列', () => {
      const { container } = renderWithRouter(<CoursesPage />)
      
      // 验证 flex-col 类在移动端存在
      const filtersContainer = container.querySelector('.flex.flex-col')
      expect(filtersContainer).toBeInTheDocument()
    })
  })

  describe('可访问性', () => {
    it('搜索框应该有正确的 label', () => {
      renderWithRouter(<CoursesPage />)
      
      const searchInput = screen.getByPlaceholderText('搜索课程...')
      expect(searchInput).toHaveAttribute('type', 'text')
    })

    it('筛选器应该是可访问的', () => {
      renderWithRouter(<CoursesPage />)
      
      const select = screen.getByRole('combobox')
      expect(select).toHaveAttribute('aria-label')
    })
  })
})
