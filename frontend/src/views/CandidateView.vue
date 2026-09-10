<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { ArrowLeft, Briefcase, Mail, MapPin } from 'lucide-vue-next'
import { api } from '../lib/api'
import { auth } from '../lib/auth'
import type { User } from '../type'
import VideoPlayer from '../components/VideoPlayer.vue'

const route = useRoute()
const id = route.params.id as string

const candidate = ref<User | null>(null)
const videoUrl = ref('')
const loading = ref(true)
const error = ref('')
const showContact = ref(false)

onMounted(async () => {
  try {
    candidate.value = await api.user(id)
    videoUrl.value = (await api.userVideo(id)).url || ''
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Profil introuvable.'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <section class="max-w-4xl mx-auto px-6 py-12">
    <RouterLink to="/candidats" class="inline-flex items-center gap-2 text-sm text-primary font-marianne font-medium hover:underline mb-8">
      <ArrowLeft class="w-4 h-4" />
      Retour
    </RouterLink>

    <p v-if="loading" class="text-sm text-text-muted font-marianne">Chargement…</p>
    <p v-else-if="error" class="alert-error">{{ error }}</p>

    <template v-else-if="candidate">
      <div class="card mb-8">
        <div class="border-b border-border">
          <VideoPlayer v-if="videoUrl" :url="videoUrl" />
          <div v-else class="w-full aspect-video bg-surface flex items-center justify-center text-text-muted font-marianne text-sm">
            Aucune vidéo de présentation
          </div>
        </div>

        <div class="p-8">
          <div class="flex items-start justify-between gap-4 flex-wrap">
            <h1 class="text-2xl font-marianne font-bold text-primary flex items-center gap-3 flex-wrap">
              {{ candidate.name }}
              <span
                class="px-2 py-0.5 text-xs uppercase tracking-wide font-bold border"
                :class="candidate.role === 'recruiter' ? 'border-primary text-primary' : 'border-action text-action'"
              >
                {{ candidate.role === 'recruiter' ? 'Recruteur' : 'Candidat' }}
              </span>
            </h1>
            <button
              v-if="!auth.isCandidate.value"
              class="btn-action text-sm shrink-0"
              @click="showContact = true"
            >
              <Mail class="w-4 h-4" />
              Contacter
            </button>
          </div>

          <div class="flex flex-col gap-2 mt-4 font-marianne text-sm text-text-muted">
            <span class="flex items-center gap-2"><MapPin class="w-4 h-4" /> {{ candidate.city || 'Non renseignée' }}</span>
            <span class="flex items-center gap-2"><Briefcase class="w-4 h-4" /> {{ candidate.sector || 'Non renseigné' }}</span>
          </div>

          <div class="mt-6 pt-6 border-t border-border">
            <h2 class="text-sm font-marianne font-bold text-primary mb-3">Compétences</h2>
            <div class="flex flex-wrap gap-2">
              <span
                v-for="(skill, i) in candidate.skills"
                :key="i"
                class="px-3 py-1 bg-surface border border-border text-primary font-marianne text-xs font-medium"
              >
                {{ skill }}
              </span>
              <p v-if="!candidate.skills?.length" class="text-text-muted font-spectral text-sm italic">
                Aucune compétence renseignée.
              </p>
            </div>
          </div>
        </div>
      </div>
    </template>

    <div
      v-if="showContact && candidate"
      class="fixed inset-0 z-50 bg-black/50 flex items-center justify-center px-4"
      @click.self="showContact = false"
    >
      <div class="card p-8 max-w-sm w-full text-center">
        <Mail class="w-10 h-10 text-action mx-auto mb-3" />
        <h2 class="text-lg font-marianne font-bold text-primary mb-1">Contacter {{ candidate.name }}</h2>
        <a
          :href="`mailto:${candidate.email}`"
          class="block text-primary font-marianne font-bold break-all mb-6"
        >
          {{ candidate.email }}
        </a>
        <div class="flex gap-3 justify-center">
          <a :href="`mailto:${candidate.email}`" class="btn-action text-sm">Envoyer un email</a>
          <button class="btn-secondary text-sm" @click="showContact = false">Fermer</button>
        </div>
      </div>
    </div>
  </section>
</template>
