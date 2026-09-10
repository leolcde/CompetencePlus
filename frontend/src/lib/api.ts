import type { AuthResponse, BadgeResult, Question, User, Video } from '../type'

const BASE = '/api'

async function request<T>(path: string, opts: RequestInit = {}): Promise<T> {
  const token = localStorage.getItem('token')
  const isForm = opts.body instanceof FormData
  const res = await fetch(BASE + path, {
    ...opts,
    headers: {
      ...(opts.body && !isForm ? { 'Content-Type': 'application/json' } : {}),
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...opts.headers,
    },
  })
  if (res.status === 401 && token) {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    location.assign('/login')
    throw new Error('Session expirée, reconnectez-vous.')
  }

  const raw = await res.text()
  let data: unknown = null
  try { data = raw ? JSON.parse(raw) : null } catch { data = raw }
  if (!res.ok) throw new Error(typeof data === 'string' && data ? data : `Erreur ${res.status}`)
  return data as T
}

const body = (v: unknown) => JSON.stringify(v)

export const api = {
  login: (email: string, password: string) => request<AuthResponse>('/login', { method: 'POST', body: body({ email, password }) }),
  register: (v: {
    name: string
    email: string
    password: string
    birthday: string
    role?: string
    city?: string
    sector?: string
    skills?: string[]
  }) => request<AuthResponse>('/register', { method: 'POST', body: body(v) }),
  me: () => request<User>('/me'),

  users: () => request<User[]>('/users'),
  user: (id: number | string) => request<User>(`/users/${id}`),

  questions: () => request<Question[]>('/questions'),
  quizStart: () => request<string>('/questionnaire/start', { method: 'POST' }),
  quizAnswer: (question_id: number, choice: string) => request<string>('/questionnaire/answer', { method: 'POST', body: body({ question_id, choice }) }),
  quizValidate: () => request<BadgeResult>('/questionnaire/validate', { method: 'POST' }),
  badge: () => request<BadgeResult>('/badge'),
  myVideo: () => request<Video>('/videos'),

  addVideoLink: (url: string) => request<Video>('/videos', { method: 'POST', body: body({ url }) }),
  addVideoFile: (file: File) => {
    const fd = new FormData()
    fd.append('file', file)
    return request<Video>('/videos', { method: 'POST', body: fd })
  },

  consent: () => request<{ active: boolean }>('/consent'),
  grantConsent: () => request('/consent', { method: 'POST' }),
  revokeConsent: () => request('/consent', { method: 'DELETE' }),
}
