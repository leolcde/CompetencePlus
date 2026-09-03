import { computed, reactive } from 'vue'

/**
 * Store d'authentification (composable maison, sans Pinia).
 * L'état (utilisateur + token JWT) est persisté dans le localStorage.
 */

const LS_KEY = 'auth_user'

export interface AuthUser {
  id: string
  name: string
  email: string
  role: string
  token: string
}

function load(): AuthUser | null {
  try {
    const raw = localStorage.getItem(LS_KEY)
    return raw ? (JSON.parse(raw) as AuthUser) : null
  } catch {
    return null
  }
}

const state = reactive<{ user: AuthUser | null }>({
  user: load(),
})

function persist() {
  try {
    if (state.user) localStorage.setItem(LS_KEY, JSON.stringify(state.user))
    else localStorage.removeItem(LS_KEY)
  } catch {
    /* stockage indisponible : on ignore */
  }
}

function setUser(user: AuthUser) {
  state.user = user
  persist()
}

/** POST /auth/login -> { token, id, name, mail, role } */
async function login(mail: string, password: string) {
  const res = await fetch('/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ mail, password }),
  })
  const body = await res.json().catch(() => ({}))
  if (!res.ok) {
    throw new Error(body.error ?? `Erreur ${res.status}`)
  }
  setUser({
    id: String(body.id ?? ''),
    name: body.name ?? '',
    email: body.mail ?? mail,
    role: body.role ?? 'candidate',
    token: body.token ?? '',
  })
}

export interface MeProfile {
  id: string
  name: string
  email: string
  role: string
  dateNaissance: string
  skills: string[]
  sector: string
  location: string
  statutCertification: string
  createdAt: string
}

/** GET /auth/me -> profil complet de l'utilisateur connecté. */
async function fetchMe(): Promise<MeProfile> {
  const res = await fetch('/auth/me', { headers: authHeader() })
  const body = await res.json().catch(() => ({}))
  if (!res.ok) {
    throw new Error(body.error ?? `Erreur ${res.status}`)
  }
  const me: MeProfile = {
    id: String(body.id ?? ''),
    name: body.identite ?? '',
    email: body.email ?? '',
    role: body.role ?? 'candidate',
    dateNaissance: body.date_naissance ?? '',
    skills: Array.isArray(body.competences) ? body.competences : [],
    sector: body.secteur ?? '',
    location: body.localisation ?? '',
    statutCertification: body.statut_certification ?? '',
    createdAt: body.created_at ?? '',
  }
  // garde le store user en phase avec le backend
  if (state.user) {
    setUser({ ...state.user, name: me.name || state.user.name, email: me.email || state.user.email, role: me.role })
  }
  return me
}

function logout() {
  state.user = null
  persist()
  // futur : fetch('/auth/logout', { method: 'POST' })
}

/** En-tête Authorization à passer aux futurs appels protégés. */
function authHeader(): Record<string, string> {
  return state.user?.token ? { Authorization: `Bearer ${state.user.token}` } : {}
}

export function useAuth() {
  return {
    user: computed(() => state.user),
    token: computed(() => state.user?.token ?? ''),
    isAuthenticated: computed(() => state.user !== null),
    login,
    fetchMe,
    logout,
    setUser,
    authHeader,
  }
}
