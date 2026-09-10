<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { RouterLink } from 'vue-router'
import { AlertTriangle, Award, CheckCircle2, MapPin, Mail, Briefcase, User, ArrowLeft } from 'lucide-vue-next'
import { auth, refreshMe } from '../lib/auth'
import { api } from '../lib/api'
import type { BadgeResult } from '../type'
import VideoPlayer from '../components/VideoPlayer.vue'

const user = auth.user

const badge = ref<BadgeResult | null>(null)
const videoUrl = ref('')
const isCandidate = computed(() => user.value?.role === 'candidate')

const roleLabel = computed(() => {
  if (user.value?.role === 'recruiter') return 'Recruteur'
  if (user.value?.role === 'candidate') return 'Candidat'
  return ''
})

const hasConsent = ref(false)
const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    await refreshMe()
    const consent = await api.consent()
    hasConsent.value = consent.active
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Erreur de chargement'
  } finally {
    loading.value = false
  }

  if (isCandidate.value) {
    try {
      badge.value = await api.badge()
    } catch {
      badge.value = null
    }
    try {
      videoUrl.value = (await api.myVideo()).url || ''
    } catch {
      videoUrl.value = ''
    }
  }
})

async function revokeConsent() {
  try {
    await api.revokeConsent()
    hasConsent.value = false
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Erreur'
  }
}

async function grantConsent() {
  try {
    await api.grantConsent()
    hasConsent.value = true
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Erreur'
  }
}
</script>

<template>
  <div class="max-w-4xl mx-auto px-6 py-12 w-full">
    <RouterLink to="/candidats" class="inline-flex items-center gap-2 text-sm text-primary font-marianne font-medium hover:underline mb-8">
      <ArrowLeft class="w-4 h-4" />
      Retour
    </RouterLink>

    <p v-if="loading" class="text-sm text-text-muted font-marianne mb-4">Chargement…</p>
    <p v-if="error" class="alert-error mb-6">{{ error }}</p>

    <div class="card p-8 mb-8">
      <h2 class="text-xl font-marianne font-bold text-primary mb-6 pb-4 border-b border-border">
        Informations personnelles
      </h2>
      <div class="flex flex-col gap-4">
        <div class="flex items-center gap-3 font-marianne text-sm">
          <User class="w-4 h-4 text-text-muted shrink-0" />
          <span class="text-text-muted w-24">Nom</span>
          <span class="text-primary font-bold flex items-center gap-2 flex-wrap">
            {{ user?.name }}
            <span
              v-if="roleLabel"
              class="px-2 py-0.5 text-xs uppercase tracking-wide font-bold border"
              :class="roleLabel === 'Recruteur' ? 'border-primary text-primary' : 'border-action text-action'"
            >
              {{ roleLabel }}
            </span>
          </span>
        </div>
        <div class="flex items-center gap-3 font-marianne text-sm">
          <Mail class="w-4 h-4 text-text-muted shrink-0" />
          <span class="text-text-muted w-24">Email</span>
          <span class="text-primary">{{ user?.email }}</span>
        </div>
        <div class="flex items-center gap-3 font-marianne text-sm">
          <MapPin class="w-4 h-4 text-text-muted shrink-0" />
          <span class="text-text-muted w-24">Ville</span>
          <span class="text-primary">{{ user?.city || "Non renseignée" }}</span>
        </div>
        <div class="flex items-center gap-3 font-marianne text-sm">
          <Briefcase class="w-4 h-4 text-text-muted shrink-0" />
          <span class="text-text-muted w-24">Secteur</span>
          <span class="text-primary">{{ user?.sector || "Non renseigné" }}</span>
        </div>
      </div>

      <div class="mt-6 pt-6 border-t border-border">
        <h3 class="text-sm font-marianne font-bold text-primary mb-3">Compétences</h3>
        <div class="flex flex-wrap gap-2">
          <span
            v-for="(skill, idx) in user?.skills"
            :key="idx"
            class="px-3 py-1 bg-surface border border-border text-primary font-marianne text-xs font-medium"
          >
            {{ skill }}
          </span>
          <p v-if="!user?.skills?.length" class="text-text-muted font-spectral text-sm italic">
            Aucune compétence renseignée.
          </p>
        </div>
      </div>
    </div>

    <!-- Badge de certification (candidat) -->
    <div v-if="isCandidate" class="card p-8 mb-8">
      <h2 class="text-xl font-marianne font-bold text-primary mb-6 pb-4 border-b border-border">
        Certification savoir-être
      </h2>

      <div v-if="badge && badge.badge" class="flex items-center gap-4">
        <Award class="w-12 h-12 text-action shrink-0" />
        <div>
          <p class="font-marianne font-bold text-primary">Badge obtenu</p>
          <p class="font-spectral text-sm text-text-muted">Score : {{ badge.score }} / 20</p>
        </div>
      </div>

      <div v-else class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <p class="font-spectral text-sm text-text-muted">
          <template v-if="badge">Badge non obtenu (score {{ badge.score }} / 20).</template>
          <template v-else>Vous n'avez pas encore passé le questionnaire.</template>
        </p>
        <RouterLink :to="{ name: 'quiz' }" class="btn-secondary text-sm shrink-0">
          {{ badge ? 'Repasser le questionnaire' : 'Passer le questionnaire' }}
        </RouterLink>
      </div>
    </div>

    <div v-if="isCandidate" class="card mb-8">
      <div class="border-b border-border">
        <VideoPlayer v-if="videoUrl" :url="videoUrl" />
        <div v-else class="w-full aspect-video bg-surface flex flex-col items-center justify-center text-text-muted p-6 text-center">
          <AlertTriangle class="w-12 h-12 mb-4" />
          <p class="font-marianne font-bold text-lg mb-2">Vidéo non disponible</p>
          <p class="font-spectral mb-4">Aucune vidéo uploadée.</p>
          <RouterLink :to="{ name: 'my-video' }" class="btn-action text-sm">Ajouter une vidéo</RouterLink>
        </div>
      </div>

      <div class="p-8">
        <h2 class="text-xl font-marianne font-bold text-primary mb-6">Ma vidéo de présentation</h2>

        <div class="flex items-start gap-4 mb-6">
          <input
            :checked="hasConsent"
            type="checkbox"
            id="consent-checkbox"
            class="mt-1 w-5 h-5 accent-primary"
            @change="hasConsent ? revokeConsent() : grantConsent()"
          />
          <div>
            <label for="consent-checkbox" class="font-marianne font-bold text-primary block mb-2 cursor-pointer">
              Consentement à la publication de la vidéo
            </label>
            <p class="font-spectral text-sm text-text-main leading-relaxed">
              J'accepte expressément que ProfilsActifs diffuse ma vidéo de présentation sur la plateforme
              à destination des recruteurs. Je consens à l'utilisation de mon image et de ma voix dans le
              cadre exclusif de la mise en relation emploi.
            </p>
          </div>
        </div>

        <div class="flex flex-col sm:flex-row sm:items-center justify-between pt-6 border-t border-border gap-4">
          <div v-if="hasConsent" class="flex items-center gap-2 text-success font-marianne font-bold text-sm">
            <CheckCircle2 class="w-5 h-5" />
            Consentement donné
          </div>
          <div v-else class="text-text-muted font-marianne text-sm italic">
            Aucun consentement actif. La vidéo est masquée.
          </div>

          <button
            :disabled="!hasConsent"
            class="btn-secondary text-sm disabled:opacity-50 disabled:cursor-not-allowed"
            @click="revokeConsent"
          >
            Révoquer mon consentement
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
