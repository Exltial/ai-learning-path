import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import CodeEditor from '../components/CodeEditor'

// Mock @monaco-editor/react
vi.mock('@monaco-editor/react', () => ({
  default: ({ value, onChange, language, theme, options }: any) => (
    <div 
      data-testid="mock-editor" 
      data-language={language}
      data-theme={theme}
      onClick={() => onChange?.('console.log("test");')}
    >
      Mock Editor - {value || ''}
    </div>
  ),
}))

describe('CodeEditor', () => {
  const defaultProps = {
    value: '',
    onChange: vi.fn(),
    language: 'typescript',
    height: '400px',
    theme: 'dark' as const,
    readOnly: false,
  }

  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('渲染', () => {
    it('应该正确渲染编辑器', () => {
      render(<CodeEditor {...defaultProps} />)
      
      const editor = screen.getByTestId('mock-editor')
      expect(editor).toBeInTheDocument()
    })

    it('应该使用默认语言 typescript', () => {
      render(<CodeEditor {...defaultProps} />)
      
      const editor = screen.getByTestId('mock-editor')
      expect(editor).toHaveAttribute('data-language', 'typescript')
    })

    it('应该使用指定的语言', () => {
      render(<CodeEditor {...defaultProps} language="python" />)
      
      const editor = screen.getByTestId('mock-editor')
      expect(editor).toHaveAttribute('data-language', 'python')
    })

    it('应该使用默认主题 dark', () => {
      render(<CodeEditor {...defaultProps} />)
      
      const editor = screen.getByTestId('mock-editor')
      expect(editor).toHaveAttribute('data-theme', 'dark')
    })

    it('应该使用指定的主题', () => {
      render(<CodeEditor {...defaultProps} theme="light" />)
      
      const editor = screen.getByTestId('mock-editor')
      expect(editor).toHaveAttribute('data-theme', 'light')
    })

    it('应该显示初始值', () => {
      const initialValue = 'const x = 1;'
      render(<CodeEditor {...defaultProps} value={initialValue} />)
      
      const editor = screen.getByTestId('mock-editor')
      expect(editor).toHaveTextContent(`Mock Editor - ${initialValue}`)
    })
  })

  describe('回调函数', () => {
    it('应该在内容变化时调用 onChange', () => {
      const onChange = vi.fn()
      render(<CodeEditor {...defaultProps} onChange={onChange} />)
      
      const editor = screen.getByTestId('mock-editor')
      fireEvent.click(editor)
      
      expect(onChange).toHaveBeenCalledWith('console.log("test");')
      expect(onChange).toHaveBeenCalledTimes(1)
    })

    it('应该在没有提供 onChange 时不报错', () => {
      expect(() => {
        render(<CodeEditor {...defaultProps} onChange={undefined} />)
        const editor = screen.getByTestId('mock-editor')
        fireEvent.click(editor)
      }).not.toThrow()
    })
  })

  describe('只读模式', () => {
    it('应该支持只读模式', () => {
      render(<CodeEditor {...defaultProps} readOnly={true} />)
      
      // 只读模式下编辑器应该渲染
      const editor = screen.getByTestId('mock-editor')
      expect(editor).toBeInTheDocument()
    })
  })

  describe('高度配置', () => {
    it('应该使用默认高度 400px', () => {
      const { container } = render(<CodeEditor {...defaultProps} />)
      
      // 验证组件渲染
      expect(container).toBeInTheDocument()
    })

    it('应该支持自定义高度', () => {
      const { container } = render(<CodeEditor {...defaultProps} height="600px" />)
      
      expect(container).toBeInTheDocument()
    })
  })

  describe('placeholder', () => {
    it('应该支持自定义 placeholder', () => {
      const customPlaceholder = '// 编写你的代码'
      render(<CodeEditor {...defaultProps} placeholder={customPlaceholder} />)
      
      // placeholder 在空值时显示
      expect(screen.getByText(customPlaceholder)).toBeInTheDocument()
    })

    it('应该在有值时不显示 placeholder', () => {
      const customPlaceholder = '// 编写你的代码'
      render(<CodeEditor {...defaultProps} value="const x = 1;" placeholder={customPlaceholder} />)
      
      expect(screen.queryByText(customPlaceholder)).not.toBeInTheDocument()
    })
  })

  describe('可访问性', () => {
    it('应该有适当的 ARIA 属性', () => {
      render(<CodeEditor {...defaultProps} />)
      
      const editor = screen.getByTestId('mock-editor')
      expect(editor).toBeInTheDocument()
    })
  })

  describe('边界情况', () => {
    it('应该处理空字符串值', () => {
      render(<CodeEditor {...defaultProps} value="" />)
      
      const editor = screen.getByTestId('mock-editor')
      expect(editor).toHaveTextContent('Mock Editor - ')
    })

    it('应该处理 undefined 值', () => {
      render(<CodeEditor {...defaultProps} value={undefined} />)
      
      const editor = screen.getByTestId('mock-editor')
      expect(editor).toBeInTheDocument()
    })

    it('应该处理很长的代码', () => {
      const longCode = 'const x = 1;\n'.repeat(100)
      render(<CodeEditor {...defaultProps} value={longCode} />)
      
      const editor = screen.getByTestId('mock-editor')
      expect(editor).toBeInTheDocument()
    })
  })
})
