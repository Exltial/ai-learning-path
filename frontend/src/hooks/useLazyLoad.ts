import { useState, useEffect, useRef, useCallback } from 'react'

interface LazyLoadOptions {
  root?: Element | null
  rootMargin?: string
  threshold?: number | number[]
  placeholder?: string
  loadOffset?: number // Load image when this many pixels before viewport
}

interface LazyLoadResult {
  ref: React.RefObject<HTMLElement>
  isVisible: boolean
  isLoaded: boolean
  error: Error | null
}

/**
 * useLazyLoad - Custom hook for lazy loading images and components
 * Uses Intersection Observer API for efficient viewport detection
 * 
 * @param options - Lazy load configuration options
 * @returns Lazy load state and ref
 * 
 * @example
 * // Lazy load an image
 * function LazyImage({ src, alt }: { src: string; alt: string }) {
 *   const { ref, isVisible, isLoaded } = useLazyLoad({ threshold: 0.1 })
 *   
 *   return (
 *     <div ref={ref} style={{ minHeight: '200px' }}>
 *       {isVisible && (
 *         <img 
 *           src={src} 
 *           alt={alt}
 *           loading="lazy"
 *           style={{ opacity: isLoaded ? 1 : 0, transition: 'opacity 0.3s' }}
 *         />
 *       )}
 *       {!isLoaded && <div className="placeholder">Loading...</div>}
 *     </div>
 *   )
 * }
 * 
 * @example
 * // Lazy load a component
 * function LazyComponent() {
 *   const { ref, isVisible } = useLazyLoad({ threshold: 0 })
 *   
 *   return (
 *     <div ref={ref}>
 *       {isVisible && <HeavyComponent />}
 *     </div>
 *   )
 * }
 */
export function useLazyLoad(options: LazyLoadOptions = {}): LazyLoadResult {
  const {
    root = null,
    rootMargin = '0px',
    threshold = 0.1,
    placeholder = '',
    loadOffset = 100,
  } = options

  const ref = useRef<HTMLElement>(null)
  const [isVisible, setIsVisible] = useState(false)
  const [isLoaded, setIsLoaded] = useState(false)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    const element = ref.current
    if (!element) return

    // Adjust threshold with load offset
    const adjustedThreshold = Array.isArray(threshold) 
      ? threshold 
      : [threshold]

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            setIsVisible(true)
            
            // Unobserve after becoming visible (load once)
            observer.unobserve(entry.target)
          }
        })
      },
      {
        root,
        rootMargin: `${loadOffset}px`,
        threshold: adjustedThreshold,
      }
    )

    observer.observe(element)

    return () => {
      observer.disconnect()
    }
  }, [root, threshold, loadOffset])

  // Mark as loaded when visible
  useEffect(() => {
    if (isVisible) {
      // Simulate loading delay for smooth UX
      const timer = setTimeout(() => {
        setIsLoaded(true)
      }, 100)
      
      return () => clearTimeout(timer)
    }
  }, [isVisible])

  return { ref, isVisible, isLoaded, error }
}

/**
 * useLazyLoadImage - Specialized hook for lazy loading images
 * Includes preloading and error handling
 * 
 * @param src - Image source URL
 * @param options - Lazy load options
 * @returns Image lazy load state
 * 
 * @example
 * function LazyImage({ src, alt, className }: { src: string; alt: string; className?: string }) {
 *   const { ref, isVisible, isLoaded, error, loadImage } = useLazyLoadImage(src)
 *   
 *   return (
 *     <div 
 *       ref={ref}
 *       className={`image-container ${className}`}
 *       style={{ minHeight: '200px', background: '#f0f0f0' }}
 *     >
 *       {isVisible && (
 *         <>
 *           {!isLoaded && !error && <div className="loading">Loading...</div>}
 *           {error && <div className="error">Failed to load</div>}
 *           <img
 *             src={src}
 *             alt={alt}
 *             loading="lazy"
 *             onLoad={loadImage}
 *             onError={() => setError(new Error('Image load failed'))}
 *             style={{
 *               display: isLoaded ? 'block' : 'none',
 *               opacity: isLoaded ? 1 : 0,
 *               transition: 'opacity 0.3s ease-in-out'
 *             }}
 *           />
 *         </>
 *       )}
 *     </div>
 *   )
 * }
 */
interface LazyLoadImageResult extends LazyLoadResult {
  loadImage: () => void
}

export function useLazyLoadImage(src: string, options: LazyLoadOptions = {}): LazyLoadImageResult {
  const lazyLoad = useLazyLoad(options)
  const { ref, isVisible, isLoaded, error } = lazyLoad
  const [imageLoaded, setImageLoaded] = useState(false)

  const loadImage = useCallback(() => {
    setImageLoaded(true)
  }, [])

  return {
    ...lazyLoad,
    isLoaded: isVisible && imageLoaded && isLoaded,
    loadImage,
  }
}

/**
 * LazyImage - Ready-to-use lazy loading image component
 * 
 * @example
 * <LazyImage
 *   src="/images/course-thumbnail.jpg"
 *   alt="Course thumbnail"
 *   className="course-image"
 *   placeholder="/images/placeholder.jpg"
 * />
 */
interface LazyImageProps extends React.ImgHTMLAttributes<HTMLImageElement> {
  src: string
  alt: string
  placeholder?: string
  rootMargin?: string
  threshold?: number
}

export function LazyImage({ 
  src, 
  alt, 
  placeholder,
  rootMargin = '100px',
  threshold = 0.1,
  className,
  style,
  ...props 
}: LazyImageProps) {
  const { ref, isVisible, isLoaded, error } = useLazyLoadImage(src, {
    rootMargin,
    threshold,
  })

  return (
    <div
      ref={ref}
      className={`lazy-image-container ${className || ''}`}
      style={{
        position: 'relative',
        minHeight: '200px',
        background: '#f5f5f5',
        ...style,
      }}
    >
      {/* Placeholder */}
      {(!isLoaded || !isVisible) && placeholder && (
        <img
          src={placeholder}
          alt=""
          className="lazy-image-placeholder"
          style={{
            position: 'absolute',
            top: 0,
            left: 0,
            width: '100%',
            height: '100%',
            objectFit: 'cover',
          }}
        />
      )}

      {/* Loading indicator */}
      {isVisible && !isLoaded && !error && (
        <div
          className="lazy-image-loading"
          style={{
            position: 'absolute',
            top: '50%',
            left: '50%',
            transform: 'translate(-50%, -50%)',
            color: '#666',
          }}
        >
          Loading...
        </div>
      )}

      {/* Error state */}
      {error && (
        <div
          className="lazy-image-error"
          style={{
            position: 'absolute',
            top: '50%',
            left: '50%',
            transform: 'translate(-50%, -50%)',
            color: '#dc3545',
            textAlign: 'center',
          }}
        >
          Failed to load image
        </div>
      )}

      {/* Actual image */}
      {isVisible && (
        <img
          src={src}
          alt={alt}
          loading="lazy"
          {...props}
          onLoad={() => {
            // Mark as loaded
            const event = new Event('load')
            props.onLoad?.(event as any)
          }}
          style={{
            width: '100%',
            height: '100%',
            objectFit: 'cover',
            opacity: isLoaded ? 1 : 0,
            transition: 'opacity 0.3s ease-in-out',
            ...style,
          }}
        />
      )}
    </div>
  )
}

/**
 * useComponentLazyLoad - Hook for lazy loading heavy components
 * Delays rendering until component is visible
 * 
 * @example
 * function HeavyComponent() {
 *   // Expensive rendering logic
 *   return <div>Heavy content</div>
 * }
 * 
 * function Page() {
 *   const { shouldRender } = useComponentLazyLoad()
 *   
 *   return (
 *     <div>
 *       <LightComponent />
 *       {shouldRender && <HeavyComponent />}
 *     </div>
 *   )
 * }
 */
interface ComponentLazyLoadOptions extends LazyLoadOptions {
  fallback?: React.ReactNode
}

interface ComponentLazyLoadResult {
  ref: React.RefObject<HTMLElement>
  shouldRender: boolean
}

export function useComponentLazyLoad(
  options: ComponentLazyLoadOptions = {}
): ComponentLazyLoadResult {
  const { ref, isVisible } = useLazyLoad(options)
  
  return {
    ref,
    shouldRender: isVisible,
  }
}

/**
 * LazyComponent - HOC for lazy loading components
 * 
 * @example
 * const LazyHeavyComponent = lazyComponent(HeavyComponent, {
 *   threshold: 0.1,
 *   fallback: <div>Loading...</div>
 * })
 * 
 * function Page() {
 *   return (
 *     <div>
 *       <LazyHeavyComponent />
 *     </div>
 *   )
 * }
 */
export function lazyComponent<P extends object>(
  Component: React.ComponentType<P>,
  options: ComponentLazyLoadOptions = {}
) {
  const { fallback, ...lazyOptions } = options

  return function LazyComponentWrapper(props: P) {
    const { ref, shouldRender } = useComponentLazyLoad(lazyOptions)

    return (
      <div ref={ref}>
        {shouldRender ? (
          <Component {...props} />
        ) : (
          fallback || <div>Loading...</div>
        )}
      </div>
    )
  }
}

export default useLazyLoad
