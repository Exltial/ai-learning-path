import React from 'react'

interface SelectProps {
  value?: string
  defaultValue?: string
  onValueChange?: (value: string) => void
  children: React.ReactNode
}

interface SelectTriggerProps {
  children: React.ReactNode
  className?: string
}

interface SelectValueProps {
  placeholder?: string
  className?: string
}

interface SelectContentProps {
  children: React.ReactNode
  className?: string
}

interface SelectItemProps {
  value: string
  children: React.ReactNode
  className?: string
}

export function Select({ value, defaultValue, onValueChange, children }: SelectProps) {
  const [selectedValue, setSelectedValue] = React.useState(defaultValue || value || '')
  
  const handleValueChange = (newValue: string) => {
    setSelectedValue(newValue)
    onValueChange?.(newValue)
  }
  
  return (
    <div className="relative" data-value={selectedValue} data-onchange={onValueChange ? 'true' : 'false'}>
      {React.Children.map(children, child => {
        if (React.isValidElement(child)) {
          return React.cloneElement(child as React.ReactElement<any>, {
            value: selectedValue,
            onValueChange: handleValueChange
          })
        }
        return child
      })}
    </div>
  )
}

export function SelectTrigger({ children, className = '' }: SelectTriggerProps) {
  return (
    <button className={`w-full inline-flex items-center justify-between px-3 py-2 border border-gray-300 rounded-lg bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 ${className}`}>
      {children}
    </button>
  )
}

export function SelectValue({ placeholder, className = '' }: SelectValueProps) {
  return (
    <span className={`text-gray-700 ${className}`}>
      {placeholder}
    </span>
  )
}

export function SelectContent({ children, className = '' }: SelectContentProps) {
  return (
    <div className={`absolute z-50 w-full mt-1 bg-white border border-gray-300 rounded-lg shadow-lg max-h-60 overflow-auto ${className}`}>
      {children}
    </div>
  )
}

export function SelectItem({ value, children, className = '' }: SelectItemProps) {
  return (
    <div 
      className={`px-3 py-2 cursor-pointer hover:bg-gray-100 ${className}`}
      data-value={value}
    >
      {children}
    </div>
  )
}
