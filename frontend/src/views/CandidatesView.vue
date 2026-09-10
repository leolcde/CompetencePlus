<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { MapPin } from 'lucide-vue-next'
import { api } from '../lib/api'
import type { User } from '../type'

const users = ref<User[]>([])
const loading = ref(true)
const error = ref('')

const PAGE_SIZE = 20
const page = ref(1)

const paginated = computed(() => {
  const start = (page.value - 1) * PAGE_SIZE
  return users.value.slice(start, start + PAGE_SIZE)
})

const totalPages = computed(() => Math.ceil(users.value.length / PAGE_SIZE))

onMounted(async () => {
  try {
    const all = await api.users()
    users.value = all.filter(u => u.role === 'candidate')
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Erreur de chargement'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="max-w-5xl mx-auto px-6 py-12 w-full">
    <h1 class="text-3xl font-marianne font-bold text-primary mb-8">Les candidats</h1>

    <p v-if="loading" class="text-sm text-text-muted font-marianne">Chargement…</p>
    <p v-if="error" class="alert-error mb-6">{{ error }}</p>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
      <RouterLink
        v-for="user in paginated"
        :key="user.id"
        :to="`/candidats/${user.id}`"
        class="card p-6 hover:shadow-md transition-shadow flex flex-col gap-3"
      >
        <div>
          <h2 class="text-lg font-marianne font-bold text-primary">{{ user.name }}</h2>
          <p class="text-sm font-marianne text-text-muted">{{ user.sector }}</p>
        </div>
        <div class="flex items-center gap-1.5 text-text-muted font-marianne text-sm">
          <MapPin class="w-4 h-4" />
          {{ user.city }}
        </div>
        <div class="flex flex-wrap gap-2 mt-2">
          <span
            v-for="skill in user.skills?.slice(0, 3)"
            :key="skill"
            class="px-2 py-1 bg-surface border border-border text-primary font-marianne text-xs"
          >
            {{ skill }}
          </span>
        </div>
      </RouterLink>
    </div>

    <p v-if="!loading && users.length === 0" class="text-text-muted font-spectral text-sm italic mt-8">
      Aucun candidat disponible pour le moment.
    </p>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="flex items-center justify-center gap-4 mt-12">
      <button
        class="btn-secondary text-sm px-4 py-2 disabled:opacity-50"
        :disabled="page === 1"
        @click="page--"
      >
        ← Précédent
      </button>
      <span class="font-marianne text-sm text-text-muted">
        Page {{ page }} / {{ totalPages }}
      </span>
      <button
        class="btn-secondary text-sm px-4 py-2 disabled:opacity-50"
        :disabled="page === totalPages"
        @click="page++"
      >
        Suivant →
      </button>
    </div>
  </div>
</template>
