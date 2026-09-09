export type Role = 'candidate' | 'recruiter' | 'admin'
export type Status = 'youth' | 'adult'

export interface User {id: number; name: string; email: string; birthday: string; status: Status; skills: string[] | null; sector: string; city: string; role: Role; created_at: string}
export interface Question { id: number; content: string; options: string[]; weight: number }
export interface BadgeResult { user_id: number; score: number; badge: boolean }
export interface Video { id: number; user_id: number; url: string; created_at: string }
export interface AuthResponse { token: string; user: User }
