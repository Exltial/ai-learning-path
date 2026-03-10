import { useEffect, useRef, useCallback } from 'react';

export interface GestureState {
  isSwiping: boolean;
  direction: 'left' | 'right' | 'up' | 'down' | null;
  startX: number;
  startY: number;
  currentX: number;
  currentY: number;
  distance: number;
}

export interface UseSwipeOptions {
  onSwipeLeft?: () => void;
  onSwipeRight?: () => void;
  onSwipeUp?: () => void;
  onSwipeDown?: () => void;
  threshold?: number;
  disabled?: boolean;
}

/**
 * 触摸手势 Hook - 支持滑动操作
 */
export function useSwipe({
  onSwipeLeft,
  onSwipeRight,
  onSwipeUp,
  onSwipeDown,
  threshold = 50,
  disabled = false,
}: UseSwipeOptions = {}) {
  const touchStartRef = useRef<GestureState>({
    isSwiping: false,
    direction: null,
    startX: 0,
    startY: 0,
    currentX: 0,
    currentY: 0,
    distance: 0,
  });

  const handleTouchStart = useCallback((e: TouchEvent) => {
    if (disabled) return;
    
    const touch = e.touches[0];
    touchStartRef.current = {
      isSwiping: true,
      direction: null,
      startX: touch.clientX,
      startY: touch.clientY,
      currentX: touch.clientX,
      currentY: touch.clientY,
      distance: 0,
    };
  }, [disabled]);

  const handleTouchMove = useCallback((e: TouchEvent) => {
    if (!touchStartRef.current.isSwiping || disabled) return;
    
    const touch = e.touches[0];
    const deltaX = touch.clientX - touchStartRef.current.startX;
    const deltaY = touch.clientY - touchStartRef.current.startY;
    
    touchStartRef.current.currentX = touch.clientX;
    touchStartRef.current.currentY = touch.clientY;
    touchStartRef.current.distance = Math.sqrt(deltaX * deltaX + deltaY * deltaY);
    
    // 确定滑动方向
    if (!touchStartRef.current.direction && touchStartRef.current.distance > 10) {
      if (Math.abs(deltaX) > Math.abs(deltaY)) {
        touchStartRef.current.direction = deltaX > 0 ? 'right' : 'left';
      } else {
        touchStartRef.current.direction = deltaY > 0 ? 'down' : 'up';
      }
    }
  }, [disabled]);

  const handleTouchEnd = useCallback(() => {
    if (!touchStartRef.current.isSwiping || disabled) return;
    
    const { direction, distance } = touchStartRef.current;
    
    if (distance >= threshold) {
      switch (direction) {
        case 'left':
          onSwipeLeft?.();
          break;
        case 'right':
          onSwipeRight?.();
          break;
        case 'up':
          onSwipeUp?.();
          break;
        case 'down':
          onSwipeDown?.();
          break;
      }
    }
    
    touchStartRef.current.isSwiping = false;
    touchStartRef.current.direction = null;
  }, [disabled, threshold, onSwipeLeft, onSwipeRight, onSwipeUp, onSwipeDown]);

  useEffect(() => {
    if (disabled) return;
    
    document.addEventListener('touchstart', handleTouchStart, { passive: true });
    document.addEventListener('touchmove', handleTouchMove, { passive: true });
    document.addEventListener('touchend', handleTouchEnd);
    
    return () => {
      document.removeEventListener('touchstart', handleTouchStart);
      document.removeEventListener('touchmove', handleTouchMove);
      document.removeEventListener('touchend', handleTouchEnd);
    };
  }, [handleTouchStart, handleTouchMove, handleTouchEnd, disabled]);

  return {
    isSwiping: touchStartRef.current.isSwiping,
    direction: touchStartRef.current.direction,
  };
}

/**
 * 双击手势 Hook
 */
export function useDoubleTap(onDoubleTap: () => void, disabled = false) {
  const lastTapRef = useRef<number>(0);
  const timerRef = useRef<number | null>(null);

  const handleTouchEnd = useCallback((e: TouchEvent) => {
    if (disabled) return;
    
    const now = Date.now();
    const deltaTime = now - lastTapRef.current;
    
    if (deltaTime < 300 && deltaTime > 0) {
      // 双击
      if (timerRef.current) {
        clearTimeout(timerRef.current);
        timerRef.current = null;
      }
      onDoubleTap();
      lastTapRef.current = 0;
    } else {
      // 单击，等待可能的第二次点击
      lastTapRef.current = now;
      timerRef.current = window.setTimeout(() => {
        lastTapRef.current = 0;
        timerRef.current = null;
      }, 300);
    }
  }, [disabled, onDoubleTap]);

  useEffect(() => {
    if (disabled) return;
    
    document.addEventListener('touchend', handleTouchEnd, { passive: true });
    
    return () => {
      document.removeEventListener('touchend', handleTouchEnd);
      if (timerRef.current) {
        clearTimeout(timerRef.current);
      }
    };
  }, [handleTouchEnd, disabled]);
}

/**
 * 捏合缩放手势 Hook
 */
export function usePinch(onPinch: (scale: number) => void, disabled = false) {
  const initialDistanceRef = useRef<number>(0);
  const currentScaleRef = useRef<number>(1);

  const getDistance = (touch1: Touch, touch2: Touch) => {
    const dx = touch1.clientX - touch2.clientX;
    const dy = touch1.clientY - touch2.clientY;
    return Math.sqrt(dx * dx + dy * dy);
  };

  const handleTouchStart = useCallback((e: TouchEvent) => {
    if (disabled || e.touches.length !== 2) return;
    
    initialDistanceRef.current = getDistance(e.touches[0], e.touches[1]);
  }, [disabled]);

  const handleTouchMove = useCallback((e: TouchEvent) => {
    if (disabled || e.touches.length !== 2) return;
    
    const currentDistance = getDistance(e.touches[0], e.touches[1]);
    const scale = currentDistance / initialDistanceRef.current;
    
    if (scale > 0) {
      currentScaleRef.current = scale;
      onPinch(scale);
    }
  }, [disabled, onPinch]);

  const handleTouchEnd = useCallback(() => {
    if (disabled) return;
    initialDistanceRef.current = 0;
  }, [disabled]);

  useEffect(() => {
    if (disabled) return;
    
    document.addEventListener('touchstart', handleTouchStart, { passive: true });
    document.addEventListener('touchmove', handleTouchMove, { passive: true });
    document.addEventListener('touchend', handleTouchEnd);
    
    return () => {
      document.removeEventListener('touchstart', handleTouchStart);
      document.removeEventListener('touchmove', handleTouchMove);
      document.removeEventListener('touchend', handleTouchEnd);
    };
  }, [handleTouchStart, handleTouchMove, handleTouchEnd, disabled]);

  return { scale: currentScaleRef.current };
}

/**
 * 长按手势 Hook
 */
export function useLongPress(
  onLongPress: () => void,
  delay = 500,
  disabled = false
) {
  const timerRef = useRef<number | null>(null);
  const isLongPressRef = useRef<boolean>(false);

  const handleTouchStart = useCallback(() => {
    if (disabled) return;
    
    isLongPressRef.current = false;
    timerRef.current = window.setTimeout(() => {
      isLongPressRef.current = true;
      onLongPress();
    }, delay);
  }, [disabled, delay, onLongPress]);

  const handleTouchEnd = useCallback(() => {
    if (disabled) return;
    
    if (timerRef.current) {
      clearTimeout(timerRef.current);
      timerRef.current = null;
    }
    isLongPressRef.current = false;
  }, [disabled]);

  const handleTouchMove = useCallback(() => {
    if (disabled) return;
    
    // 如果手指移动，取消长按
    if (timerRef.current) {
      clearTimeout(timerRef.current);
      timerRef.current = null;
    }
    isLongPressRef.current = false;
  }, [disabled]);

  useEffect(() => {
    if (disabled) return;
    
    document.addEventListener('touchstart', handleTouchStart, { passive: true });
    document.addEventListener('touchend', handleTouchEnd);
    document.addEventListener('touchmove', handleTouchMove, { passive: true });
    
    return () => {
      document.removeEventListener('touchstart', handleTouchStart);
      document.removeEventListener('touchend', handleTouchEnd);
      document.removeEventListener('touchmove', handleTouchMove);
      if (timerRef.current) {
        clearTimeout(timerRef.current);
      }
    };
  }, [handleTouchStart, handleTouchEnd, handleTouchMove, disabled]);

  return { isLongPressing: isLongPressRef.current };
}
