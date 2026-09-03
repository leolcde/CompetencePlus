<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import {
  CheckCircle2,
  AlertTriangle,
  Eye,
  ShieldCheck,
  Video,
} from 'lucide-vue-next'
import { MOCK_PROFILES } from '../assets/data/mock'
import JebBadge from '../assets/JebBadge.vue'

// Profil connecté (démo : le profil "1")
const source = MOCK_PROFILES.find((p) => p.id === '1') ?? MOCK_PROFILES[0]

const form = reactive({
  name: source.name,
  job: source.job,
  city: source.city,
  skills: source.skills.join(', '),
})

const hasConsent = ref(source.hasConsent)
const saved = ref(false)
const profileViews = 128

const skillList = computed(() =>
  form.skills
    .split(',')
    .map((s) => s.trim())
    .filter(Boolean),
)

const inputClass =
  'w-full border border-border p-3 rounded-none focus:outline-none focus:border-primary font-spectral'
const labelClass = 'block font-marianne font-bold text-primary mb-2 text-sm'

function save() {
  saved.value = true
  setTimeout(() => (saved.value = false), 2500)
}
</script>

<template>
  <div class="flex-1 bg-surface py-10 sm:py-16 px-6">
    <div class="max-w-5xl mx-auto space-y-10">
      <!-- En-tête -->
      <div class="flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
          <p class="font-marianne font-bold text-action text-xs uppercase tracking-wide mb-2">
            Espace personnel
          </p>
          <h1 class="text-3xl sm:text-4xl font-marianne font-black text-primary tracking-tight">
            {{ form.name }}
          </h1>
          <p class="text-text-muted font-spectral mt-1">{{ form.job }} · {{ form.city }}</p>
        </div>
        <RouterLink :to="`/profil/${source.id}`" class="btn-secondary text-sm">
          Voir mon profil public
        </RouterLink>
      </div>

      <!-- Statistiques -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div class="bg-white border border-border p-6">
          <div class="flex items-center gap-2 text-text-muted font-marianne text-sm mb-3">
            <ShieldCheck class="w-4 h-4" />
            Certification JEB
          </div>
          <p class="text-3xl font-marianne font-black text-primary">
            {{ source.isCertified ? `${source.score}/100` : 'Non passée' }}
          </p>
        </div>

        <div class="bg-white border border-border p-6">
          <div class="flex items-center gap-2 text-text-muted font-marianne text-sm mb-3">
            <Video class="w-4 h-4" />
            Statut de la vidéo
          </div>
          <p
            class="text-3xl font-marianne font-black"
            :class="hasConsent ? 'text-success' : 'text-text-muted'"
          >
            {{ hasConsent ? 'Publiée' : 'Masquée' }}
          </p>
        </div>

        <div class="bg-white border border-border p-6">
          <div class="flex items-center gap-2 text-text-muted font-marianne text-sm mb-3">
            <Eye class="w-4 h-4" />
            Vues du profil
          </div>
          <p class="text-3xl font-marianne font-black text-primary">{{ profileViews }}</p>
        </div>
      </div>

      <!-- Certification -->
      <div
        v-if="source.isCertified"
        class="bg-white border border-border p-8 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4"
      >
        <div class="flex items-center gap-4">
          <JebBadge :large="true" />
          <p class="font-spectral text-text-main">
            Votre savoir-être professionnel est certifié. Score de {{ source.score }}/100.
          </p>
        </div>
      </div>
      <div v-else class="bg-white border border-border p-8">
        <h2 class="text-xl font-marianne font-bold text-primary mb-2">Passez la certification JEB</h2>
        <p class="font-spectral text-text-main mb-6">
          Valorisez vos compétences douces auprès des recruteurs.
        </p>
        <button class="btn-action text-sm">Commencer le test</button>
      </div>

      <!-- Informations du profil -->
      <div class="bg-white border border-border p-8 md:p-10">
        <h2 class="text-2xl font-marianne font-bold text-primary mb-8">Informations du profil</h2>

        <p
          v-if="saved"
          class="mb-6 border border-success bg-surface text-success font-marianne text-sm p-3 flex items-center gap-2"
        >
          <CheckCircle2 class="w-4 h-4" />
          Modifications enregistrées.
        </p>

        <form class="space-y-6" @submit.prevent="save">
          <div>
            <label :class="labelClass">Nom complet</label>
            <input v-model="form.name" type="text" :class="inputClass" required />
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label :class="labelClass">Intitulé de poste</label>
              <input v-model="form.job" type="text" :class="inputClass" />
            </div>
            <div>
              <label :class="labelClass">Ville</label>
              <input v-model="form.city" type="text" :class="inputClass" />
            </div>
          </div>

          <div>
            <label :class="labelClass">Compétences</label>
            <input
              v-model="form.skills"
              type="text"
              :class="inputClass"
              placeholder="Communication, Logistique, Management"
            />
            <div class="flex flex-wrap gap-2 mt-3">
              <span
                v-for="(skill, idx) in skillList"
                :key="idx"
                class="px-3 py-1 bg-surface border border-border text-primary font-marianne text-xs font-medium"
              >
                {{ skill }}
              </span>
            </div>
          </div>

          <button type="submit" class="btn-action mt-4">Enregistrer les modifications</button>
        </form>
      </div>

      <!-- Gestion de la vidéo / consentement -->
      <div class="bg-white border border-border p-8 md:p-10">
        <h2 class="text-2xl font-marianne font-bold text-primary mb-2">Ma vidéo de présentation</h2>
        <p class="font-spectral text-text-main mb-8">
          Vous contrôlez la diffusion de votre vidéo auprès des recruteurs.
        </p>

        <div class="aspect-video bg-surface border border-border flex items-center justify-center mb-8">
          <img
            v-if="hasConsent"
            :src="source.videoUrl"
            :alt="`Vidéo de ${form.name}`"
            class="w-full h-full object-cover grayscale opacity-90"
          />
          <div v-else class="flex flex-col items-center text-text-muted p-6 text-center">
            <AlertTriangle class="w-10 h-10 mb-3" />
            <p class="font-marianne font-bold">Vidéo masquée</p>
          </div>
        </div>

        <div class="border border-border p-6 bg-surface">
          <label class="flex items-start gap-3 cursor-pointer">
            <input
              v-model="hasConsent"
              type="checkbox"
              class="mt-1 w-5 h-5 border-border rounded-none text-action focus:ring-action"
            />
            <span class="font-spectral text-sm text-text-main leading-relaxed">
              J'accepte expressément que ProfilsActifs diffuse ma vidéo de présentation sur la
              plateforme à destination des recruteurs, et l'utilisation de mon image et de ma voix
              dans le cadre exclusif de la mise en relation emploi.
            </span>
          </label>

          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pt-6 mt-6 border-t border-border">
            <div
              v-if="hasConsent"
              class="flex items-center gap-2 text-success font-marianne font-bold text-sm"
            >
              <CheckCircle2 class="w-5 h-5" />
              Consentement actif — vidéo publiée
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
  </div>
</template>
