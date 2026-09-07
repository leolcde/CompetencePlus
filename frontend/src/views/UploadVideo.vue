<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { CheckCircle2, Info, Link2, Upload, Video, X } from 'lucide-vue-next'
import { useAuth } from '../stores/auth'

const { can, authHeader } = useAuth()

const MAX_MB = 100
const ACCEPT = 'video/mp4,video/webm,video/quicktime,video/x-matroska'

const mode = ref<'file' | 'link'>('file')
const consent = ref(false)
const loading = ref(false)
const error = ref('')
const done = ref(false)

const url = ref('')
const trimmed = computed(() => url.value.trim())
const isValidUrl = computed(() => {
  if (!trimmed.value) return false
  try {
    const u = new URL(trimmed.value)
    return u.protocol === 'http:' || u.protocol === 'https:'
  } catch {
    return false
  }
})
const showPreview = ref(false)
const embed = computed<{ kind: 'youtube' | 'vimeo' | 'file' | 'unknown'; src: string }>(() => {
  if (!isValidUrl.value) return { kind: 'unknown', src: '' }
  const u = new URL(trimmed.value)
  const host = u.hostname.replace(/^www\./, '')
  if (host === 'youtube.com' || host === 'm.youtube.com') {
    const id = u.searchParams.get('v')
    if (id) return { kind: 'youtube', src: `https://www.youtube.com/embed/${id}` }
  }
  if (host === 'youtu.be') {
    const id = u.pathname.slice(1)
    if (id) return { kind: 'youtube', src: `https://www.youtube.com/embed/${id}` }
  }
  if (host === 'vimeo.com') {
    const id = u.pathname.split('/').filter(Boolean)[0]
    if (id && /^\d+$/.test(id)) return { kind: 'vimeo', src: `https://player.vimeo.com/video/${id}` }
  }
  if (/\.(mp4|webm|ogg|mov)$/i.test(u.pathname)) return { kind: 'file', src: trimmed.value }
  return { kind: 'unknown', src: trimmed.value }
})

const file = ref<File | null>(null)
const objectUrl = ref('')
const dragOver = ref(false)

watch(file, (f) => {
  if (objectUrl.value) URL.revokeObjectURL(objectUrl.value)
  objectUrl.value = f ? URL.createObjectURL(f) : ''
})
onBeforeUnmount(() => {
  if (objectUrl.value) URL.revokeObjectURL(objectUrl.value)
})

function pickFile(f: File | undefined | null) {
  error.value = ''
  if (!f) return
  if (!f.type.startsWith('video/')) {
    error.value = 'Le fichier doit être une vidéo (mp4, webm, mov…).'
    return
  }
  if (f.size > MAX_MB * 1024 * 1024) {
    error.value = `La vidéo dépasse ${MAX_MB} Mo.`
    return
  }
  file.value = f
}
function onInputChange(e: Event) {
  pickFile((e.target as HTMLInputElement).files?.[0])
}
function onDrop(e: DragEvent) {
  dragOver.value = false
  pickFile(e.dataTransfer?.files?.[0])
}
function clearFile() {
  file.value = null
}
function humanSize(bytes: number) {
  return bytes < 1024 * 1024
    ? `${Math.round(bytes / 1024)} Ko`
    : `${(bytes / 1024 / 1024).toFixed(1)} Mo`
}

const canSubmit = computed(() => {
  if (!consent.value || loading.value) return false
  return mode.value === 'file' ? !!file.value : isValidUrl.value
})

async function submit() {
  error.value = ''
  if (!canSubmit.value) {
    error.value = 'Complétez le formulaire et acceptez le consentement.'
    return
  }

  loading.value = true
  try {
    let res: Response
    if (mode.value === 'file' && file.value) {
      const fd = new FormData()
      fd.append('video', file.value)
      fd.append('consent', 'true')
      res = await fetch('/profils/video', { method: 'POST', headers: { ...authHeader() }, body: fd })
    } else {
      res = await fetch('/profils/video', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', ...authHeader() },
        body: JSON.stringify({ url: trimmed.value, consent: true }),
      })
    }

    if (res.ok) {
      done.value = true
      return
    }
    if (res.status === 404 || res.status === 501) {
      error.value =
        "L'enregistrement des vidéos sera bientôt disponible. Vos informations n'ont pas été perdues."
      return
    }
    const body = await res.json().catch(() => ({}))
    error.value = body.error ?? `Erreur ${res.status}`
  } catch {
    error.value = 'Impossible de contacter le serveur. Réessayez plus tard.'
  } finally {
    loading.value = false
  }
}

function reset() {
  done.value = false
  showPreview.value = false
  url.value = ''
  file.value = null
  consent.value = false
}
</script>

<template>
  <div class="flex-1 bg-surface py-10 sm:py-16 px-6">
    <div class="max-w-2xl mx-auto">
      <RouterLink
        :to="{ name: 'dashboard' }"
        class="text-sm text-primary font-marianne font-medium hover:underline"
      >
        ← Mon espace
      </RouterLink>

      <h1 class="text-3xl sm:text-4xl font-marianne font-black text-primary tracking-tight mt-3">
        Ma vidéo de présentation
      </h1>
      <p class="font-spectral text-text-muted mt-2">
        Une vidéo courte (1 minute) qui vous présente au-delà du CV.
      </p>

      <div v-if="!can.publishVideo" class="card p-8 mt-8 text-center">
        <Info class="w-8 h-8 text-secondary mx-auto mb-3" />
        <p class="font-marianne font-bold text-primary">Réservé aux candidats</p>
        <p class="font-spectral text-text-muted text-sm mt-1">
          Seuls les comptes candidat peuvent publier une vidéo de présentation.
        </p>
      </div>

      <div v-else-if="done" class="card p-8 mt-8 text-center">
        <CheckCircle2 class="w-10 h-10 text-success mx-auto mb-3" />
        <p class="font-marianne font-bold text-primary text-lg">Vidéo enregistrée</p>
        <p class="font-spectral text-text-muted text-sm mt-1">
          Votre vidéo est désormais visible par les recruteurs sur votre profil.
        </p>
        <div class="flex flex-col sm:flex-row gap-3 justify-center mt-6">
          <RouterLink :to="{ name: 'dashboard' }" class="btn-action">Retour à mon espace</RouterLink>
          <button type="button" class="btn-secondary" @click="reset">Remplacer la vidéo</button>
        </div>
      </div>

      <form v-else class="card p-6 sm:p-8 mt-8 space-y-6" @submit.prevent="submit">
        <p v-if="error" class="alert-error">{{ error }}</p>

        <div class="inline-flex border border-border rounded-md overflow-hidden font-marianne text-sm">
          <button
            type="button"
            class="px-4 py-2 transition-colors"
            :class="mode === 'file' ? 'bg-primary text-white' : 'bg-white text-primary hover:bg-surface'"
            @click="mode = 'file'"
          >
            Depuis mon ordinateur
          </button>
          <button
            type="button"
            class="px-4 py-2 border-l border-border transition-colors"
            :class="mode === 'link' ? 'bg-primary text-white' : 'bg-white text-primary hover:bg-surface'"
            @click="mode = 'link'"
          >
            Coller un lien
          </button>
        </div>

        <div v-if="mode === 'file'">
          <label
            class="field-label"
            for="video-file"
          >Fichier vidéo</label>

          <label
            v-if="!file"
            for="video-file"
            class="flex flex-col items-center justify-center gap-2 border-2 border-dashed rounded-md py-10 px-4 text-center cursor-pointer transition-colors"
            :class="dragOver ? 'border-primary bg-primary/5' : 'border-border hover:border-primary/40 hover:bg-surface'"
            @dragover.prevent="dragOver = true"
            @dragleave.prevent="dragOver = false"
            @drop.prevent="onDrop"
          >
            <Upload class="w-7 h-7 text-secondary" />
            <span class="font-marianne font-bold text-primary text-sm">
              Cliquez pour choisir un fichier
            </span>
            <span class="font-spectral text-xs text-text-muted">
              ou glissez-le ici — mp4, webm, mov · {{ MAX_MB }} Mo max
            </span>
          </label>

          <input
            id="video-file"
            type="file"
            :accept="ACCEPT"
            class="sr-only"
            @change="onInputChange"
          />

          <div v-if="file" class="space-y-3">
            <div class="flex items-center justify-between gap-3 border border-border rounded-md px-3 py-2">
              <div class="min-w-0">
                <p class="font-marianne text-sm text-primary truncate">{{ file.name }}</p>
                <p class="font-spectral text-xs text-text-muted">{{ humanSize(file.size) }}</p>
              </div>
              <button
                type="button"
                class="shrink-0 text-text-muted hover:text-danger"
                aria-label="Retirer le fichier"
                @click="clearFile"
              >
                <X class="w-4 h-4" />
              </button>
            </div>
            <video
              :src="objectUrl"
              controls
              class="w-full aspect-video bg-black rounded-md border border-border"
            />
          </div>
        </div>

        <div v-else>
          <label for="video-url" class="field-label">Lien de la vidéo</label>
          <div class="relative">
            <Link2 class="w-4 h-4 text-text-muted absolute left-3 top-1/2 -translate-y-1/2" />
            <input
              id="video-url"
              v-model="url"
              type="url"
              inputmode="url"
              placeholder="https://www.youtube.com/watch?v=…"
              class="field pl-9"
              :class="trimmed && !isValidUrl ? 'border-danger focus:border-danger focus:ring-danger/15' : ''"
            />
          </div>
          <p class="text-xs text-text-muted font-spectral mt-1.5">YouTube, Vimeo ou lien direct .mp4 / .webm.</p>
          <p v-if="trimmed && !isValidUrl" class="text-xs text-danger font-marianne mt-1">Lien invalide.</p>

          <div v-if="isValidUrl" class="mt-3">
            <button
              type="button"
              class="text-sm font-marianne font-medium text-action hover:underline inline-flex items-center gap-1.5"
              @click="showPreview = !showPreview"
            >
              <Video class="w-4 h-4" />
              {{ showPreview ? "Masquer l'aperçu" : 'Prévisualiser' }}
            </button>
            <div v-if="showPreview" class="mt-3 aspect-video bg-black rounded-md overflow-hidden border border-border">
              <iframe
                v-if="embed.kind === 'youtube' || embed.kind === 'vimeo'"
                :src="embed.src"
                class="w-full h-full"
                allow="fullscreen"
                referrerpolicy="strict-origin-when-cross-origin"
              />
              <video v-else-if="embed.kind === 'file'" :src="embed.src" controls class="w-full h-full" />
              <div
                v-else
                class="w-full h-full flex items-center justify-center text-white/70 text-sm font-marianne p-4 text-center"
              >
                Aperçu indisponible pour ce lien — il sera tout de même enregistré.
              </div>
            </div>
          </div>
        </div>

        <div class="pt-5 border-t border-border">
          <label class="flex items-start gap-3 cursor-pointer">
            <input v-model="consent" type="checkbox" class="mt-1 w-5 h-5 accent-primary" />
            <span class="font-spectral text-sm text-text-main leading-relaxed">
              J'accepte expressément que ProfilsActifs diffuse ma vidéo de présentation sur la
              plateforme à destination des recruteurs, et l'utilisation de mon image et de ma voix
              dans le cadre exclusif de la mise en relation emploi.
            </span>
          </label>
        </div>

        <button type="submit" class="btn-action w-full" :disabled="!canSubmit">
          {{ loading ? 'Enregistrement…' : 'Publier ma vidéo' }}
        </button>
      </form>
    </div>
  </div>
</template>
