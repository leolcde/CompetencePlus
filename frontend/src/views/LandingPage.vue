<script setup lang="ts">
import Logo from '../components/Logo.vue'
import { useAuth } from '../stores/auth'

const { isAuthenticated } = useAuth()

const logoSrc = '/logo.png'
const year = new Date().getFullYear()

const footerLinks = [
  { label: 'Accueil', to: { name: 'home' } },
  { label: 'Les profils', to: { name: 'profiles' } },
  { label: "S'inscrire", to: { name: 'signup' } },
  { label: 'Mon espace', to: { name: 'dashboard' } },
]
</script>

<template>
  <div class="flex-1 flex flex-col">
    <section class="relative flex-1 flex flex-col items-center justify-center px-6 py-20 text-center overflow-hidden">
      <img
        :src="logoSrc"
        alt=""
        aria-hidden="true"
        class="pointer-events-none select-none absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 w-[min(90vw,520px)] opacity-[0.06]"
      />
      <div class="relative z-10 w-full max-w-2xl mx-auto flex flex-col items-center">
        <h1 class="text-4xl sm:text-5xl font-black text-primary font-marianne tracking-tight mb-6 leading-tight">
          La mise en relation<br />simple et transparente.
        </h1>
        <p class="text-lg sm:text-xl text-text-main font-spectral leading-relaxed mb-10">
          ProfilsActifs est la plateforme gouvernementale permettant aux demandeurs d'emploi de se
          présenter en vidéo et aux recruteurs de découvrir des talents authentiques, certifiés JEB.
        </p>
        <div class="flex flex-col sm:flex-row gap-4 justify-center w-full max-w-md">
          <RouterLink
            v-if="!isAuthenticated"
            to="/signup"
            class="btn-action w-full text-lg py-3.5"
          >
            Inscription
          </RouterLink>
          <RouterLink
            v-else
            to="/dashboard"
            class="btn-action w-full text-lg py-3.5"
          >
            Mon espace
          </RouterLink>
          <RouterLink to="/profils" class="btn-secondary w-full text-lg py-3.5">
            Découvrir les profils
          </RouterLink>
        </div>
      </div>
    </section>

    <footer class="bg-surface border-t border-border">
      <div class="max-w-7xl mx-auto px-6 py-8 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <Logo :size="28" with-wordmark />

        <nav class="flex flex-wrap gap-x-6 gap-y-2">
          <RouterLink
            v-for="link in footerLinks"
            :key="link.label"
            :to="link.to"
            class="text-text-muted hover:text-primary font-marianne text-sm hover:underline underline-offset-4 transition-colors"
          >
            {{ link.label }}
          </RouterLink>
        </nav>

        <p class="text-text-muted font-marianne text-xs">
          © {{ year }} ProfilsActifs — République française
        </p>
      </div>
    </footer>
  </div>
</template>
