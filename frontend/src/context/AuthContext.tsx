import { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import type { User, AuthResponse } from '../types'
import api from '../api/client'

interface AuthContextValue {
  user: User | null
  token: string | null
  login: (email: string, password: string) => Promise<void>
  register: (data: RegisterData) => Promise<void>
  logout: () => void
  isLoading: boolean
}

interface RegisterData {
  full_name: string
  email: string
  phone: string
  password: string
  role: string
  city: string
  country: string
}

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(localStorage.getItem('skroda_token'))
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    if (token) {
      api.get<User>('/me')
        .then((res) => setUser(res.data))
        .catch(() => logout())
        .finally(() => setIsLoading(false))
    } else {
      setIsLoading(false)
    }
  }, [token])

  const login = async (email: string, password: string) => {
    const res = await api.post<AuthResponse>('/auth/login', { email, password })
    localStorage.setItem('skroda_token', res.data.token)
    setToken(res.data.token)
    setUser(res.data.user)
  }

  const register = async (data: RegisterData) => {
    const res = await api.post<AuthResponse>('/auth/register', data)
    localStorage.setItem('skroda_token', res.data.token)
    setToken(res.data.token)
    setUser(res.data.user)
  }

  const logout = () => {
    localStorage.removeItem('skroda_token')
    setToken(null)
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ user, token, login, register, logout, isLoading }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
