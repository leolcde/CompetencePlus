<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { ArrowLeft, MapPin, MessageSquare, AlertTriangle, CheckCircle2 } from 'lucide-vue-next'
import { MOCK_PROFILES } from '../assets/data/mock'
import JebBadge from '../assets/JebBadge.vue'

const route = useRoute()
const id = String(route.params.id)
const profile = MOCK_PROFILES.find(p => p.id === id) ?? MOCK_PROFILES[0]
const isMyProfile = id === '1'
const hasConsent = ref(profile.hasConsent)
</script>

<template>
  <div class="max-w-4xl mx-auto px-6 py-12 w-full">
    <RouterLink to="/feed" class="inline-flex items-center gap-2 text-sm text-primary font-marianne font-medium hover:underline mb-8">
      <ArrowLeft class="w-4 h-4" />
      Retour au feed
    </RouterLink>

    <div class="bg-white border border-border">
      <!-- Zone vidéo -->
      <div class="w-full aspect-video bg-surface relative border-b border-border flex items-center justify-center">
        <img
          v-if="hasConsent"
          :src="profile.videoUrl"
          :alt="`Vidéo de présentation de ${profile.name}`"
          class="w-full h-full object-cover grayscale opacity-90"
        />
        <div v-else class="flex flex-col items-center justify-center text-text-muted p-6 text-center">
          <AlertTriangle class="w-12 h-12 mb-4" />
          <p class="font-marianne font-bold text-lg mb-2">Vidéo non disponible</p>
          <p class="font-spectral">Le consentement de publication a été révoqué.</p>
        </div>

        <div v-if="hasConsent" class="absolute inset-0 flex items-center justify-center">
          <button class="w-20 h-20 bg-action text-action-foreground rounded-full flex items-center justify-center pl-1.5 hover:scale-105 transition-transform shadow-lg">
            <svg class="w-8 h-8" fill="currentColor" viewBox="0 0 20 20"><path d="M4 4l12 6-12 6z" /></svg>
          </button>
        </div>
      </div>

      <div class="p-8 md:p-12">
        <div class="flex flex-col md:flex-row justify-between items-start gap-8 mb-12">
          <div>
            <h1 class="text-3xl sm:text-4xl font-marianne font-black text-primary mb-2 tracking-tight">{{ profile.name }}</h1>
            <p class="text-xl font-marianne text-text-main font-medium mb-6">{{ profile.job }}</p>
            <div class="flex items-center gap-2 text-text-muted font-marianne text-sm">
              <MapPin class="w-4 h-4" />
              {{ profile.city }}
            </div>
          </div>

          <div class="flex flex-col gap-4 min-w-[200px] shrink-0">
            <button class="btn-action w-full flex items-center justify-center gap-2">
              <MessageSquare class="w-4 h-4" />
              Contacter
            </button>

            <div v-if="profile.isCertified" class="w-full p-4 border border-success/30 bg-success/5 flex flex-col items-center gap-2">
              <JebBadge :large="true" />
              <span class="font-marianne font-bold text-success text-sm mt-1">
                Score : {{ profile.score }}/100
              </span>
            </div>
          </div>
        </div>

        <div class="pt-8 border-t border-border">
          <h2 class="text-2xl font-marianne font-bold text-primary mb-6">Compétences</h2>
          <div class="flex flex-wrap gap-2">
            <span
              v-for="(skill, idx) in profile.skills"
              :key="idx"
              class="px-4 py-2 bg-surface border border-border text-primary font-marianne text-sm font-medium"
            >
              {{ skill }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Zone de gestion (candidat uniquement) -->
    <div v-if="isMyProfile" class="mt-12 p-8 border border-border bg-surface">
      <h2 class="text-xl font-marianne font-bold text-primary mb-6">Gestion de ma vidéo (Zone Privée)</h2>

      <div class="bg-white p-6 border border-border">
        <div class="flex items-start gap-4 mb-6">
          <input
            v-model="hasConsent"
            type="checkbox"
            id="consent-checkbox"
            class="mt-1 w-5 h-5 border-border rounded-none text-action focus:ring-action"
          />
          <div>
            <label for="consent-checkbox" class="font-marianne font-bold text-primary block mb-2 cursor-pointer">
              Consentement à la publication de la vidéo
            </label>
            <p class="font-spectral text-sm text-text-main leading-relaxed">
              J'accepte expressément que ProfilsActifs diffuse ma vidéo de présentation sur la plateforme à destination des recruteurs. Je consens à l'utilisation de mon image et de ma voix dans le cadre exclusif de la mise en relation emploi.
            </p>
          </div>
        </div>

        <div class="flex flex-col sm:flex-row sm:items-center justify-between pt-6 border-t border-border gap-4">
          <div v-if="hasConsent" class="flex items-center gap-2 text-success font-marianne font-bold text-sm">
            <CheckCircle2 class="w-5 h-5" />
            Consentement donné le {{ profile.consentDate ?? 'récemment' }}
          </div>
          <div v-else class="text-text-muted font-marianne text-sm italic">
            Aucun consentement actif. La vidéo est masquée.
          </div>

          <button
            :disabled="!hasConsent"
            class="btn-secondary text-sm disabled:opacity-50 disabled:cursor-not-allowed"
            @click="hasConsent = false"
          >
            Révoquer mon consentement
          </button>
        </div>
      </div>
    </div>
  </div>
</template>