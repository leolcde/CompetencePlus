<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { ArrowLeft, Briefcase, MapPin } from 'lucide-vue-next'
import { api } from '../lib/api'
import type { User } from '../type'
import VideoPlayer from '../components/VideoPlayer.vue'

const route = useRoute()
const id = route.params.id as string

const candidate = ref<User | null>(null)
const videoUrl = ref('')
const loading = ref(true)
const error = ref('')

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
          <h1 class="text-2xl font-marianne font-bold text-primary flex items-center gap-3 flex-wrap">
            {{ candidate.name }}
            <span
              class="px-2 py-0.5 text-xs uppercase tracking-wide font-bold border"
              :class="candidate.role === 'recruiter' ? 'border-primary text-primary' : 'border-action text-action'"
            >
              {{ candidate.role === 'recruiter' ? 'Recruteur' : 'Candidat' }}
            </span>
          </h1>

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
  </section>
</template>
