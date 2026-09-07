<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import {
  CheckCircle2,
  AlertTriangle,
  Eye,
  MessageSquare,
  ShieldCheck,
  Video,
} from 'lucide-vue-next'
import { MOCK_PROFILES } from '../assets/data/mock'
import { useAuth } from '../stores/auth'
import { getCert } from '../stores/certification'
import JebBadge from '../assets/JebBadge.vue'

const { user, fetchMe, isRecruiter, can } = useAuth()

const cert = getCert(user.value?.id)

// Champs non encore fournis par le backend : fallback sur le mock.
const source = MOCK_PROFILES.find((p) => p.id === user.value?.id) ?? MOCK_PROFILES[0]

const form = reactive({
  id: user.value?.id ?? source.id,
  name: user.value?.name || source.name,
  job: source.job,
  city: source.city,
  skills: source.skills.join(', '),
})

const hasConsent = ref(false)
const saved = ref(false)
const profileViews = 128
const contactedCount = 0
const loadError = ref('')

onMounted(async () => {
  try {
    const me = await fetchMe()
    form.id = me.id || form.id
    form.name = me.name || form.name
    form.job = me.sector || form.job
    form.city = me.location || form.city
    if (me.skills.length) form.skills = me.skills.join(', ')
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : 'Impossible de charger votre profil'
  }
})

const skillList = computed(() =>
  form.skills
    .split(',')
    .map((s) => s.trim())
    .filter(Boolean),
)

const inputClass =
  'field'
const labelClass = 'field-label'

function save() {
  saved.value = true
  setTimeout(() => (saved.value = false), 2500)
}
</script>

<template>
  <div class="flex-1 bg-surface py-10 sm:py-16 px-6">
    <div class="max-w-5xl mx-auto space-y-10">
      <div class="flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
          <p class="font-marianne font-bold text-action text-xs uppercase tracking-wide mb-2">
            Paramètres du profil
          </p>
          <h1 class="text-3xl sm:text-4xl font-marianne font-black text-primary tracking-tight flex items-center gap-3 flex-wrap">
            {{ form.name }}
            <span
              class="px-2 py-0.5 text-xs uppercase tracking-wide font-bold border"
              :class="isRecruiter ? 'border-primary text-primary' : 'border-action text-action'"
            >
              {{ isRecruiter ? 'Recruteur' : 'Candidat' }}
            </span>
          </h1>
          <p class="text-text-muted font-spectral mt-1">{{ form.job }} · {{ form.city }}</p>
        </div>
        <RouterLink :to="`/profil/${form.id}`" class="btn-secondary text-sm">
          Voir mon profil public
        </RouterLink>
      </div>

      <p
        v-if="loadError"
        class="alert-error"
      >
        {{ loadError }}
      </p>

      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <template v-if="isRecruiter">
          <div class="card p-6">
            <div class="flex items-center gap-2 text-text-muted font-marianne text-sm mb-3">
              <MessageSquare class="w-4 h-4" />
              Candidats contactés
            </div>
            <p class="text-3xl font-marianne font-black text-primary">{{ contactedCount }}</p>
          </div>
        </template>
        <template v-else>
          <div class="card p-6">
            <div class="flex items-center gap-2 text-text-muted font-marianne text-sm mb-3">
              <ShieldCheck class="w-4 h-4" />
              Certification JEB
            </div>
            <p class="text-3xl font-marianne font-black text-primary">
              {{ cert?.done ? (cert.badge ? 'Certifié' : 'Non certifié') : 'Non passée' }}
            </p>
            <p v-if="cert?.done" class="font-marianne text-xs text-text-muted mt-1">
              Score : {{ cert.score }}
            </p>
          </div>

          <div class="card p-6">
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
        </template>

        <div class="card p-6">
          <div class="flex items-center gap-2 text-text-muted font-marianne text-sm mb-3">
            <Eye class="w-4 h-4" />
            Vues du profil
          </div>
          <p class="text-3xl font-marianne font-black text-primary">{{ profileViews }}</p>
        </div>
      </div>

      <div v-if="isRecruiter" class="card p-8">
        <h2 class="text-xl font-marianne font-bold text-primary mb-2">Espace recruteur</h2>
        <p class="font-spectral text-text-main mb-6">
          Parcourez les profils pour découvrir des candidats, contactez-les et likez ceux qui
          vous intéressent. Vous ne publiez pas de vidéo de présentation.
        </p>
        <RouterLink :to="{ name: 'profiles' }" class="btn-action text-sm">Voir les profils</RouterLink>
      </div>

      <div
        v-if="!isRecruiter && cert?.done && cert.badge"
        class="card p-8 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4"
      >
        <div class="flex items-center gap-4">
          <JebBadge :large="true" />
          <p class="font-spectral text-text-main">
            Votre savoir-être professionnel est certifié. Score de {{ cert.score }}.
          </p>
        </div>
      </div>
      <div v-else-if="!isRecruiter && cert?.done" class="card p-8">
        <h2 class="text-xl font-marianne font-bold text-primary mb-2">Certification non obtenue</h2>
        <p class="font-spectral text-text-main mb-6">
          Score de {{ cert.score }}. Vous pouvez repasser le questionnaire pour l'améliorer.
        </p>
        <RouterLink :to="{ name: 'quiz' }" class="btn-secondary text-sm">Repasser le test</RouterLink>
      </div>
      <div v-else-if="!isRecruiter" class="card p-8">
        <h2 class="text-xl font-marianne font-bold text-primary mb-2">Passez la certification JEB</h2>
        <p class="font-spectral text-text-main mb-6">
          Valorisez vos compétences douces auprès des recruteurs.
        </p>
        <RouterLink :to="{ name: 'quiz' }" class="btn-action text-sm">Commencer le test</RouterLink>
      </div>

      <div class="card p-8 md:p-10">
        <h2 class="text-2xl font-marianne font-bold text-primary mb-8">Informations du profil</h2>

        <p
          v-if="saved"
          class="mb-6 alert-success flex items-center gap-2"
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

      <div v-if="can.publishVideo" class="card p-8 md:p-10">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
          <div>
            <h2 class="text-2xl font-marianne font-bold text-primary">Ma vidéo de présentation</h2>
            <p class="font-spectral text-text-main mt-1">
              Vous contrôlez la diffusion de votre vidéo auprès des recruteurs.
            </p>
          </div>
          <RouterLink :to="{ name: 'upload' }" class="btn-action text-sm shrink-0">
            <Video class="w-4 h-4" />
            {{ hasConsent ? 'Modifier ma vidéo' : 'Ajouter ma vidéo' }}
          </RouterLink>
        </div>

        <div class="aspect-video bg-surface border border-border rounded-md flex items-center justify-center">
          <img
            v-if="hasConsent"
            :src="source.videoUrl"
            :alt="`Vidéo de ${form.name}`"
            class="w-full h-full object-cover"
          />
          <div v-else class="flex flex-col items-center text-text-muted p-6 text-center">
            <AlertTriangle class="w-10 h-10 mb-3" />
            <p class="font-marianne font-bold">Aucune vidéo publiée</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
