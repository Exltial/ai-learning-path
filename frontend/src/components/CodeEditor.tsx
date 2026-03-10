import { useEffect, useState } from 'react'
import Editor from '@monaco-editor/react'

interface CodeEditorProps {
  value?: string
  onChange?: (value: string) => void
  language?: string
  height?: string
  theme?: 'light' | 'dark'
  readOnly?: boolean
  placeholder?: string
}

export default function CodeEditor({
  value = '',
  onChange,
  language = 'typescript',
  height = '400px',
  theme = 'dark',
  readOnly = false,
  placeholder = '// 在此编写代码...',
}: CodeEditorProps) {
  const [editorValue, setEditorValue] = useState(value)

  useEffect(() => {
    setEditorValue(value)
  }, [value])

  const handleEditorChange = (newValue: string | undefined) => {
    const val = newValue || ''
    setEditorValue(val)
    onChange?.(val)
  }

  return (
    <div className="border border-secondary-300 dark:border-secondary-700 rounded-lg overflow-hidden">
      <Editor
        height={height}
        language={language}
        value={editorValue}
        onChange={handleEditorChange}
        theme={theme === 'dark' ? 'vs-dark' : 'light'}
        options={{
          readOnly,
          minimap: { enabled: false },
          fontSize: 14,
          lineNumbers: 'on',
          scrollBeyondLastLine: false,
          automaticLayout: true,
          padding: { top: 16 },
        }}
      />
      {editorValue === '' && !readOnly && (
        <div className="absolute pointer-events-none text-secondary-400 px-4 py-2">
          {placeholder}
        </div>
      )}
    </div>
  )
}
