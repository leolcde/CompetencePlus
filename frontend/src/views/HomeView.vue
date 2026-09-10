<script setup lang="ts">
import { auth } from '../lib/auth'
</script>

<template>
    <section class="h-[calc(100dvh-4rem)] sm:h-[calc(100dvh-5rem)] snap-start relative flex flex-col items-center justify-center px-6 sm:px-12 text-center overflow-hidden">
      <div
        class="absolute inset-0 z-0 pointer-events-none opacity-[0.18]"
        style="background-image: url('https://images.unsplash.com/photo-1477959858617-67f85cf4f1df?w=2000&fit=crop&auto=format'); background-size: cover; background-position: center;"
        aria-hidden="true"
      />
      <div class="absolute bottom-0 left-0 right-0 h-32 z-10 pointer-events-none" style="background: linear-gradient(to bottom, transparent, white);" aria-hidden="true" />
      <div class="relative z-10 w-full max-w-4xl md:max-w-5xl xl:max-w-6xl 2xl:max-w-7xl mx-auto flex flex-col items-center">
        <h1 class="text-4xl sm:text-5xl md:text-6xl xl:text-7xl 2xl:text-8xl font-black text-primary font-marianne tracking-tight mb-8 leading-tight drop-shadow-sm">
          La mise en relation <br /> simple et transparente.
        </h1>
        <p class="text-xl md:text-2xl xl:text-3xl text-text-main font-spectral max-w-3xl xl:max-w-4xl mx-auto leading-relaxed font-medium mb-8 sm:mb-14">
          Compétences+ est la plateforme gouvernementale permettant aux demandeurs d'emploi de se présenter en vidéo et aux recruteurs de découvrir des talents authentiques.
        </p>
        <div class="flex flex-col sm:flex-row gap-4 sm:gap-6 justify-center w-full max-w-lg xl:max-w-xl 2xl:max-w-2xl">
          <RouterLink
            v-if="!auth.isLogged.value"
            to="/signup"
            class="btn-action w-full flex justify-center items-center text-lg shadow-md py-4"
          >
            Inscription
          </RouterLink>
          <RouterLink to="/candidats" class="btn-secondary w-full text-lg shadow-sm bg-white/90 backdrop-blur-sm py-4">
            Découvrir les profils
          </RouterLink>
        </div>
      </div>
    </section>
</template>
