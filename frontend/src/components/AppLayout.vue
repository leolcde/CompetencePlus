<script setup lang="ts">
import { RouterLink, useRouter } from 'vue-router'
import { auth, logout } from '../lib/auth'
import LegalBanner from './LegalBanner.vue'

const router = useRouter()

function onLogout() {
  logout()
  router.push({ name: 'home' })
}
</script>

<template>
  <div class="min-h-screen flex flex-col bg-surface">
    <header class="sticky top-0 z-40 bg-white border-b border-border">
      <div class="max-w-5xl mx-auto px-4 h-14 flex items-center justify-between">
        <RouterLink :to="{ name: 'home' }" class="font-marianne font-black text-primary text-lg">
          Competences+
        </RouterLink>
        <nav class="flex items-center gap-4 font-marianne text-sm">
          <RouterLink :to="{ name: 'candidates' }" class="text-primary hover:underline">Candidats</RouterLink>
          <template v-if="auth.isLogged.value">
            <RouterLink :to="{ name: 'account' }" class="text-primary hover:underline">Mon espace</RouterLink>
            <button class="btn-secondary text-xs px-3 py-1.5" @click="onLogout">Se déconnecter</button>
          </template>
          <RouterLink v-else :to="{ name: 'login' }" class="btn-action text-xs px-3 py-1.5">Se connecter</RouterLink>
        </nav>
      </div>
      <LegalBanner />
    </header>

    <main class="flex-1">
      <slot />
    </main>

    <footer class="bg-white border-t border-border">
      <div class="max-w-5xl mx-auto px-4 py-4 text-xs text-text-muted font-marianne">
        © {{ new Date().getFullYear() }} Competences+  République française
      </div>
    </footer>
  </div>
</template>
