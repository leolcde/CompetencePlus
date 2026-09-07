<script setup lang="ts">
import Logo from '../components/Logo.vue'
</script>

<template>
  <div class="flex-1 bg-surface flex items-center justify-center px-6 py-20">
    <div class="max-w-md w-full text-center">
      <Logo :size="64" class="justify-center mb-8" />

      <p class="font-marianne font-black text-primary text-6xl sm:text-7xl tracking-tight">404</p>
      <h1 class="font-marianne font-bold text-primary text-xl mt-4">Page introuvable</h1>
      <p class="font-spectral text-text-muted mt-2">
        Cette page n'existe pas ou a été déplacée.
      </p>

      <div class="flex flex-col sm:flex-row gap-3 justify-center mt-8">
        <RouterLink :to="{ name: 'home' }" class="btn-action">Retour à l'accueil</RouterLink>
        <RouterLink :to="{ name: 'profiles' }" class="btn-secondary">Voir les profils</RouterLink>
      </div>
    </div>
  </div>
</template>
