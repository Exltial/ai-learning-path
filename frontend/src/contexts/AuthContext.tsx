import { createContext, useContext, useState, useEffect, type ReactNode } from 'react'
import { api } from '@/services/api'
import type { User, LoginData, RegisterData } from '@/types'

interface AuthContextType {
  user: User | null
  loading: boolean
  login: (data: LoginData) => Promise<{ success: boolean; error?: string }>
  register: (data: RegisterData) => Promise<{ success: boolean; error?: string }>
  logout: () => void
  isAuthenticated: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const initAuth = async () => {
      const token = localStorage.getItem('auth_token')
      if (token) {
        const response = await api.getCurrentUser()
        if (response.success && response.data) {
          setUser(response.data)
        } else {
          localStorage.removeItem('auth_token')
        }
      }
      setLoading(false)
    }
    initAuth()
  }, [])

  const login = async (data: LoginData): Promise<{ success: boolean; error?: string }> => {
    const response = await api.login(data)
    if (response.success && response.data) {
      setUser(response.data.user)
      return { success: true }
    }
    return { success: false, error: response.error || '登录失败' }
  }

  const register = async (data: RegisterData): Promise<{ success: boolean; error?: string }> => {
    const response = await api.register(data)
    if (response.success && response.data) {
      setUser(response.data.user)
      return { success: true }
    }
    return { success: false, error: response.error || '注册失败' }
  }

  const logout = () => {
    api.logout()
    setUser(null)
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        login,
        register,
        logout,
        isAuthenticated: !!user,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
