import React from 'react'

interface AvatarProps {
  src?: string
  alt?: string
  fallback?: string
  size?: 'sm' | 'md' | 'lg'
  className?: string
}

export function Avatar({ 
  src, 
  alt = 'Avatar', 
  fallback = '?', 
  size = 'md',
  className = '' 
}: AvatarProps) {
  const sizeStyles = {
    sm: 'w-8 h-8 text-sm',
    md: 'w-10 h-10 text-base',
    lg: 'w-12 h-12 text-lg',
  }
  
  const [hasError, setHasError] = React.useState(false)
  
  return (
    <div 
      className={`inline-flex items-center justify-center rounded-full bg-gray-200 overflow-hidden ${sizeStyles[size]} ${className}`}
    >
      {src && !hasError ? (
        <img 
          src={src} 
          alt={alt} 
          className="w-full h-full object-cover"
          onError={() => setHasError(true)}
        />
      ) : (
        <span className="font-medium text-gray-600">{fallback}</span>
      )}
    </div>
  )
}
