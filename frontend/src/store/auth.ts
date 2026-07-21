import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export type AuthUser = {
  id: string
  tenantId: string
  email: string
  displayName: string
  role: string
}

type AuthState = {
  token: string | null
  user: AuthUser | null
  tenantCode: string | null
  setSession: (token: string, user: AuthUser, tenantCode: string) => void
  clear: () => void
}

const AUTH_STORAGE_KEY = 'metartls-auth'

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      user: null,
      tenantCode: null,
      setSession: (token, user, tenantCode) => set({ token, user, tenantCode }),
      clear: () => set({ token: null, user: null, tenantCode: null }),
    }),
    { name: AUTH_STORAGE_KEY },
  ),
)
