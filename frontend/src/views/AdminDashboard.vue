<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Users, Briefcase, Award, Video, LogOut } from 'lucide-vue-next'
import { api } from '../lib/api'
import { logout } from '../lib/auth'
import type { User } from '../type'

const router = useRouter()
const users = ref<User[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    users.value = await api.users()
  } catch {
    users.value = []
  } finally {
    loading.value = false
  }
})

const candidates = computed(() => users.value.filter(u => u.role === 'candidate'))
const recruiters = computed(() => users.value.filter(u => u.role === 'recruiter'))

function onLogout() {
  logout()
  router.push({ name: 'admin-login' })
}
</script>

<template>
  <div class="min-h-screen bg-surface">
    <header class="bg-white border-b border-border sticky top-0 z-40">
      <div class="max-w-6xl mx-auto px-6 h-14 flex items-center justify-between">
        <span class="font-marianne font-black text-primary text-lg">Administration</span>
        <button class="btn-secondary text-xs px-3 py-1.5 flex items-center gap-2" @click="onLogout">
          <LogOut class="w-3.5 h-3.5" />
          Se déconnecter
        </button>
      </div>
    </header>

    <div class="max-w-6xl mx-auto px-6 py-12">
      <h1 class="text-3xl font-marianne font-bold text-primary mb-10">Tableau de bord</h1>

      <p v-if="loading" class="text-sm text-text-muted font-marianne mb-8">Chargement…</p>

      <div class="grid grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
        <div class="card p-6 flex flex-col gap-2">
          <Users class="w-8 h-8 text-action mb-1" />
          <p class="text-3xl font-marianne font-black text-primary">{{ candidates.length }}</p>
          <p class="text-sm font-marianne text-text-muted">Candidats actifs</p>
        </div>
        <div class="card p-6 flex flex-col gap-2">
          <Briefcase class="w-8 h-8 text-primary mb-1" />
          <p class="text-3xl font-marianne font-black text-primary">{{ recruiters.length }}</p>
          <p class="text-sm font-marianne text-text-muted">Recruteurs</p>
        </div>
        <div class="card p-6 flex flex-col gap-2">
          <Award class="w-8 h-8 text-success mb-1" />
          <p class="text-3xl font-marianne font-black text-primary">—</p>
          <p class="text-sm font-marianne text-text-muted">Certifiés JEB</p>
        </div>
        <div class="card p-6 flex flex-col gap-2">
          <Video class="w-8 h-8 text-text-muted mb-1" />
          <p class="text-3xl font-marianne font-black text-primary">—</p>
          <p class="text-sm font-marianne text-text-muted">Vidéos en ligne</p>
        </div>
      </div>

      <div class="card">
        <div class="p-6 border-b border-border">
          <h2 class="text-lg font-marianne font-bold text-primary">Profils récents</h2>
        </div>
        <div class="divide-y divide-border">
          <div
            v-for="user in users.slice(0, 10)"
            :key="user.id"
            class="px-6 py-4 flex items-center justify-between"
          >
            <div>
              <p class="font-marianne font-bold text-primary text-sm">{{ user.name }}</p>
              <p class="font-marianne text-text-muted text-xs">{{ user.email }}</p>
            </div>
            <span
              class="px-2 py-0.5 text-xs uppercase tracking-wide font-bold border font-marianne"
              :class="user.role === 'candidate' ? 'border-action text-action' : user.role === 'recruiter' ? 'border-primary text-primary' : 'border-success text-success'"
            >
              {{ user.role === 'candidate' ? 'Candidat' : user.role === 'recruiter' ? 'Recruteur' : 'Admin' }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
