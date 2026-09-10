import { computed, reactive } from 'vue'
import type { User } from '../type'
import { api } from './api'

const state = reactive<{ user: User | null }>({
  user: JSON.parse(localStorage.getItem('user') || 'null'),
})

function save(token: string, user: User) {
  state.user = user
  localStorage.setItem('token', token)
  localStorage.setItem('user', JSON.stringify(user))
}

export function logout() {
  state.user = null
  localStorage.removeItem('token')
  localStorage.removeItem('user')
}

export async function login(email: string, password: string) {
  const r = await api.login(email, password)
  save(r.token, r.user)
}

export async function register(v: {
  name: string
  email: string
  password: string
  birthday: string
  role?: string
  city?: string
  sector?: string
  skills?: string[]
}) {
  const r = await api.register(v)
  save(r.token, r.user)
}

export async function refreshMe() {
  const u = await api.me()
  state.user = u
  localStorage.setItem('user', JSON.stringify(u))
}

export const auth = {
  user: computed(() => state.user),
  isLogged: computed(() => state.user !== null),
  isCandidate: computed(() => state.user?.role === 'candidate'),
}
