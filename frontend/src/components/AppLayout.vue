<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { auth, logout } from '../lib/auth'
import LegalBanner from './LegalBanner.vue'

const router = useRouter()

function onLogout() {
  logout()
  router.push({ name: 'home' })
}

const year = new Date().getFullYear()

const footerColumns = computed(() => {
  const columns = [
    {
      title: 'Navigation',
      links: [
        { label: 'Accueil', to: { name: 'home' } },
        { label: 'Les candidats', to: { name: 'candidates' } },
      ],
    },
  ]

  if (auth.isLogged.value) {
    columns.push({
      title: 'Mon compte',
      links: [
        { label: 'Mon espace', to: { name: 'account' } },
        { label: 'Le questionnaire', to: { name: 'quiz' } },
        { label: 'Ma vidéo', to: { name: 'my-video' } },
      ],
    })
  } else {
    columns.push({
      title: 'Accès',
      links: [
        { label: 'Se connecter', to: { name: 'login' } },
        { label: "S'inscrire", to: { name: 'signup' } },
      ],
    })
  }

  return columns
})
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
          <template v-else>
            <RouterLink :to="{ name: 'signup' }" class="btn-secondary text-xs px-3 py-1.5">S'inscrire</RouterLink>
            <RouterLink :to="{ name: 'login' }" class="btn-action text-xs px-3 py-1.5">Se connecter</RouterLink>
          </template>
        </nav>
      </div>

    </header>

    <main class="flex-1">
      <slot />
    </main>

    <footer class="bg-white border-t border-border mt-auto">
      <LegalBanner />
      <div class="max-w-5xl mx-auto px-6 py-10">
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-8">
          <div>
            <span class="text-lg font-marianne font-black text-primary tracking-tight">Competences+</span>
            <p class="text-text-muted font-spectral text-sm leading-relaxed mt-2">
              Plateforme de mise en relation entre demandeurs d'emploi et recruteurs :
              présentation vidéo et questionnaire de savoir-être.
            </p>
          </div>

          <div v-for="col in footerColumns" :key="col.title">
            <h3 class="font-marianne font-bold text-primary text-sm uppercase tracking-wide mb-3">
              {{ col.title }}
            </h3>
            <ul class="space-y-2">
              <li v-for="link in col.links" :key="link.label">
                <RouterLink
                  :to="link.to"
                  class="text-text-muted hover:text-primary font-marianne text-sm hover:underline underline-offset-4 transition-colors"
                >
                  {{ link.label }}
                </RouterLink>
              </li>
            </ul>
          </div>
        </div>

        <p class="mt-8 pt-6 border-t border-border text-text-muted font-marianne text-xs">
          © {{ year }} Competences+
        </p>
      </div>
    </footer>
  </div>
</template>
