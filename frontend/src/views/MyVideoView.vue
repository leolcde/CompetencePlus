<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { api } from '../lib/api'
import type { Video } from '../type'
import VideoPlayer from '../components/VideoPlayer.vue'

const mode = ref<'link' | 'file'>('link')
const url = ref('')
const file = ref<File | null>(null)

const consent = ref(false)
const loading = ref(true)
const saving = ref(false)
const error = ref('')
const saved = ref<Video | null>(null)

const missing = computed(() => {
  if (!consent.value) return 'Vous devez accepter le consentement.'
  if (mode.value === 'link' && !url.value.trim()) return 'Collez le lien de votre vidéo.'
  if (mode.value === 'file' && !file.value) return 'Choisissez un fichier vidéo.'
  return ''
})

onMounted(async () => {
  try {
    const c = await api.consent()
    consent.value = c.active
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Erreur de chargement'
  } finally {
    loading.value = false
  }
})

const consentBusy = ref(false)

async function toggleConsent() {
  if (consentBusy.value) return
  consentBusy.value = true
  error.value = ''
  try {
    if (consent.value) await api.revokeConsent()
    else await api.grantConsent()
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Erreur'
  } finally {
    try {
      consent.value = (await api.consent()).active
    } catch {
      /* garde l'état courant */
    }
    consentBusy.value = false
  }
}

function onFile(e: Event) {
  file.value = (e.target as HTMLInputElement).files?.[0] ?? null
}

async function submit() {
  error.value = ''
  saving.value = true
  try {
    if (mode.value === 'link') {
      saved.value = await api.addVideoLink(url.value.trim())
    } else if (file.value) {
      saved.value = await api.addVideoFile(file.value)
    }
    url.value = ''
    file.value = null
  } catch (e) {
    error.value = e instanceof Error ? e.message : "Échec de l'envoi"
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <section class="max-w-xl mx-auto px-4 py-12">
    <h1 class="text-2xl font-marianne font-bold mb-6">Ma vidéo de présentation</h1>

    <p v-if="loading" class="text-text-muted font-marianne">Chargement…</p>
    <p v-else-if="error" class="alert-error mb-6">{{ error }}</p>

    <template v-if="!loading">
      <div class="card p-6 mb-6">
        <label class="flex items-start gap-3 cursor-pointer">
          <input type="checkbox" class="mt-1 w-5 h-5 accent-primary" :checked="consent" :disabled="consentBusy" @change="toggleConsent" />
          <span class="font-spectral text-sm text-text-main leading-relaxed">
            J'accepte que ma vidéo de présentation soit diffusée sur la plateforme à destination des
            recruteurs, ainsi que l'utilisation de mon image et de ma voix dans le cadre de la mise
            en relation emploi.
          </span>
        </label>
      </div>

      <form class="card p-6 space-y-4" @submit.prevent="submit">
        <div class="flex gap-2">
          <button
            type="button"
            class="btn-secondary text-sm flex-1"
            :class="{ 'bg-surface': mode === 'link' }"
            @click="mode = 'link'"
          >
            Par lien
          </button>
          <button
            type="button"
            class="btn-secondary text-sm flex-1"
            :class="{ 'bg-surface': mode === 'file' }"
            @click="mode = 'file'"
          >
            Par fichier
          </button>
        </div>

        <div v-if="mode === 'link'">
          <label class="field-label">Lien de la vidéo (YouTube, Vimeo…)</label>
          <input v-model="url" type="url" class="field" placeholder="https://…" required />
        </div>

        <div v-else>
          <label class="field-label">Fichier vidéo</label>
          <input type="file" accept="video/mp4,video/webm,video/quicktime" class="field" @change="onFile" />
          <p class="text-xs text-text-muted mt-1 font-spectral">mp4, webm ou mov — 100 Mo max.</p>
        </div>

        <button type="submit" class="btn-action w-full" :disabled="saving || !!missing">
          {{ saving ? 'Envoi…' : 'Envoyer ma vidéo' }}
        </button>
        <p v-if="missing" class="text-xs text-text-muted font-spectral text-center">{{ missing }}</p>
      </form>

      <div v-if="saved" class="mt-6">
        <p class="alert-success mb-3">Vidéo enregistrée.</p>
        <VideoPlayer :url="saved.url" />
      </div>
    </template>
  </section>
</template>
