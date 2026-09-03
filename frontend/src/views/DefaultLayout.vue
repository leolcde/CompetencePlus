<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { LogOut, Menu, X } from 'lucide-vue-next'
import { useAuth } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const { user, isAuthenticated, logout } = useAuth()
const open = ref(false)
const loginDropdown = ref(false)

function closeLoginDropdown() {
  loginDropdown.value = false
}

watch(
  () => route.fullPath,
  () => { loginDropdown.value = false }
)

// referme le menu à chaque changement de page
watch(
  () => route.fullPath,
  () => {
    open.value = false
  },
)

// liens visibles selon l'état de connexion
const links = computed(() => {
  const base = [{ label: 'Feed', to: { name: 'feed' } }]
  if (isAuthenticated.value) {
    base.push({ label: 'Mon espace', to: { name: 'dashboard' } })
  } else {
    base.push({ label: "S'inscrire", to: { name: 'signup' } })
  }
  return base
})

function onLogout() {
  logout()
  router.push({ name: 'home' })
}
</script>

<template>
  <div class="min-h-screen bg-background flex flex-col">
    <header class="border-b border-border bg-white sticky top-0 z-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 h-16 sm:h-20 flex items-center justify-between">
        <RouterLink to="/" class="flex items-center shrink-0">
          <span class="text-xl sm:text-2xl font-black text-primary font-marianne tracking-tight">
            ProfilsActifs
          </span>
        </RouterLink>

        <!-- Nav desktop -->
        <nav class="hidden md:flex items-center gap-6 lg:gap-8 font-marianne font-medium text-sm">
          <RouterLink
            v-for="link in links"
            :key="link.label"
            :to="link.to"
            class="text-primary hover:underline underline-offset-4"
          >
            {{ link.label }}
          </RouterLink>

          <!-- Dropdown Connexion (non connecté) -->
          <div v-if="!isAuthenticated" class="relative">
            <button
              class="btn-action text-xs px-4 py-2 inline-flex items-center gap-1.5"
              @click="loginDropdown = !loginDropdown"
            >
              Connexion
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </button>
            <div
              v-if="loginDropdown"
              class="absolute right-0 top-full mt-2 w-48 bg-white border border-border shadow-md z-50"
              @mouseleave="closeLoginDropdown"
            >
              <RouterLink
                :to="{ name: 'login', query: { type: 'candidat' } }"
                class="block px-4 py-3 font-marianne text-sm text-primary hover:bg-surface"
              >
                Candidat
              </RouterLink>
              <RouterLink
                :to="{ name: 'login', query: { type: 'recruteur' } }"
                class="block px-4 py-3 font-marianne text-sm text-primary hover:bg-surface border-t border-border"
              >
                Recruteur
              </RouterLink>
            </div>
          </div>

          <template v-if="isAuthenticated">
            <span class="text-text-muted">·</span>
            <span class="text-text-muted">{{ user?.name || user?.email }}</span>
            <button
              class="btn-secondary text-xs px-3 py-1 border border-border inline-flex items-center gap-1.5"
              @click="onLogout"
            >
              <LogOut class="w-3.5 h-3.5" />
              Se déconnecter
            </button>
          </template>
        </nav>

        <!-- Bouton menu mobile -->
        <button
          class="md:hidden inline-flex items-center justify-center w-10 h-10 -mr-2 text-primary"
          :aria-expanded="open"
          aria-label="Menu"
          @click="open = !open"
        >
          <X v-if="open" class="w-6 h-6" />
          <Menu v-else class="w-6 h-6" />
        </button>
      </div>

      <!-- Menu mobile déroulant -->
      <nav
        v-if="open"
        class="md:hidden border-t border-border bg-white px-4 py-4 flex flex-col gap-1 font-marianne font-medium"
      >
        <RouterLink
          v-for="link in links"
          :key="link.label"
          :to="link.to"
          class="py-3 px-2 text-primary hover:bg-surface rounded"
        >
          {{ link.label }}
        </RouterLink>

        <button
          v-if="isAuthenticated"
          class="btn-secondary text-sm mt-2 justify-center inline-flex items-center gap-2"
          @click="onLogout"
        >
          <LogOut class="w-4 h-4" />
          Se déconnecter ({{ user?.name || user?.email }})
        </button>
      </nav>
    </header>

    <main class="flex-1 flex flex-col">
      <RouterView />
    </main>
  </div>
</template>
