<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Heart, Lock, MapPin, MessageSquare, Play, VideoOff } from 'lucide-vue-next'
import { MOCK_PROFILES } from '../assets/data/mock'
import JebBadge from '../assets/JebBadge.vue'
import { useAuth } from '../stores/auth'

const PER_PAGE = 20

const router = useRouter()
const route = useRoute()
const { isAuthenticated, can } = useAuth()

const showAuthModal = ref(false)

function requireAuth(): boolean {
  if (isAuthenticated.value) return true
  showAuthModal.value = true
  return false
}

// connecté mais rôle sans la permission -> action ignorée
const canLike = computed(() => !isAuthenticated.value || can.value.like)
const canContact = computed(() => !isAuthenticated.value || can.value.contact)

function chooseType(type: 'candidat' | 'recruteur') {
  showAuthModal.value = false
  router.push({ name: 'login', query: { type, redirect: route.fullPath } })
}

interface ProfileCard {
  id: string
  name: string
  job: string
  city: string
  skills: string[]
  isCertified: boolean
  score: number | null
  videoUrl: string
  hasConsent: boolean
}

// Liste mock affichée immédiatement, remplacée par GET /profils au montage.
const mockList: ProfileCard[] = MOCK_PROFILES.map((p) => ({
  id: p.id,
  name: p.name,
  job: p.job,
  city: p.city,
  skills: p.skills,
  isCertified: p.isCertified,
  score: p.score,
  videoUrl: '',
  hasConsent: false,
}))

const profiles = ref<ProfileCard[]>(mockList)

onMounted(async () => {
  try {
    const res = await fetch('/profils')
    if (!res.ok) return
    const rows = await res.json()
    if (!Array.isArray(rows) || rows.length === 0) return
    profiles.value = rows.map((r: Record<string, unknown>): ProfileCard => ({
      id: String(r.id ?? ''),
      name: (r.name as string) || 'Profil',
      job: (r.job as string) || 'Profil ProfilsActifs',
      city: (r.city as string) || 'France',
      skills: Array.isArray(r.skills) ? (r.skills as string[]) : [],
      isCertified: Boolean(r.isCertified),
      score: typeof r.score === 'number' ? r.score : null,
      videoUrl: (r.videoUrl as string) || '',
      hasConsent: Boolean(r.hasConsent) && Boolean(r.videoUrl),
    }))
  } catch {
    /* on garde la liste mock */
  }
})

// --- Pagination (client, ?page=N dans l'URL pour partage) ---
const totalPages = computed(() => Math.max(1, Math.ceil(profiles.value.length / PER_PAGE)))

const page = computed(() => {
  const n = Number(route.query.page)
  if (!Number.isInteger(n) || n < 1) return 1
  return Math.min(n, totalPages.value)
})

const pageItems = computed(() => {
  const start = (page.value - 1) * PER_PAGE
  return profiles.value.slice(start, start + PER_PAGE)
})

function goToPage(n: number) {
  if (n < 1 || n > totalPages.value || n === page.value) return
  router.push({ name: 'profiles', query: n === 1 ? {} : { page: String(n) } })
}

// si le nombre de profils rétrécit (mock -> API), on borne la page courante
watch(totalPages, (max) => {
  if (page.value > max) goToPage(max)
})

// --- Vidéo : lecture uniquement au clic (pas d'autoplay) ---
const playing = reactive<Set<string>>(new Set())
function play(id: string) {
  playing.add(id)
}

// --- Like (état local par utilisateur, pas de compteur) ---
const LS_KEY = 'feed_likes'

function loadLiked(): Set<string> {
  try {
    return new Set<string>(JSON.parse(localStorage.getItem(LS_KEY) ?? '[]'))
  } catch {
    return new Set<string>()
  }
}

const liked = reactive(loadLiked())

function toggleLike(id: string) {
  if (!requireAuth()) return
  if (!can.value.like) return
  if (liked.has(id)) liked.delete(id)
  else liked.add(id)
  try {
    localStorage.setItem(LS_KEY, JSON.stringify([...liked]))
  } catch {
    /* stockage indisponible : on ignore */
  }
}
</script>

<template>
  <div class="flex-1 bg-surface py-10 sm:py-14 px-4 sm:px-6">
    <div class="max-w-7xl mx-auto">
      <div class="flex flex-col sm:flex-row sm:items-end justify-between gap-3 mb-8">
        <div>
          <h1 class="text-3xl sm:text-4xl font-marianne font-black text-primary tracking-tight">
            Les profils
          </h1>
          <p class="text-text-muted font-spectral mt-1">
            {{ profiles.length }} profils · page {{ page }} / {{ totalPages }}
          </p>
        </div>

        <RouterLink
          v-if="!isAuthenticated"
          :to="{ name: 'login', query: { redirect: route.fullPath } }"
          class="inline-flex items-center gap-2 text-sm text-primary font-marianne font-medium border border-border bg-white px-3 py-2 hover:bg-gray-100 transition-colors"
        >
          <Lock class="w-3.5 h-3.5" />
          Connectez-vous pour aimer et contacter
        </RouterLink>
      </div>

      <!-- Grille -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5">
        <article
          v-for="item in pageItems"
          :key="item.id"
          class="bg-white border border-border flex flex-col"
        >
          <!-- Média : lecture au clic uniquement -->
          <div class="relative aspect-video bg-surface border-b border-border">
            <template v-if="item.hasConsent">
              <img
                :src="item.videoUrl"
                :alt="`Présentation de ${item.name}`"
                class="absolute inset-0 w-full h-full object-cover"
                :class="playing.has(item.id) ? '' : 'grayscale opacity-90'"
              />
              <button
                v-if="!playing.has(item.id)"
                type="button"
                class="absolute inset-0 flex items-center justify-center bg-black/20 hover:bg-black/30 transition-colors"
                :aria-label="`Lire la vidéo de ${item.name}`"
                @click="play(item.id)"
              >
                <span class="w-14 h-14 bg-action text-action-foreground rounded-full flex items-center justify-center pl-1 shadow-lg">
                  <Play class="w-6 h-6" fill="currentColor" />
                </span>
              </button>
            </template>
            <div
              v-else
              class="absolute inset-0 flex flex-col items-center justify-center text-text-muted"
            >
              <VideoOff class="w-8 h-8 mb-2" />
              <p class="font-marianne font-bold text-sm">Vidéo non disponible</p>
            </div>
          </div>

          <div class="p-4 flex flex-col flex-1">
            <div v-if="item.isCertified" class="mb-2">
              <JebBadge />
            </div>
            <h2 class="text-lg font-marianne font-bold text-primary leading-tight">
              {{ item.name }}
            </h2>
            <p class="font-marianne text-sm text-text-main">{{ item.job }}</p>
            <p class="font-marianne text-xs text-text-muted flex items-center gap-1 mt-0.5">
              <MapPin class="w-3 h-3" />
              {{ item.city }}
            </p>

            <div class="flex flex-wrap gap-1.5 mt-3">
              <span
                v-for="skill in item.skills"
                :key="skill"
                class="px-2 py-0.5 bg-surface border border-border text-primary font-marianne text-xs"
              >
                {{ skill }}
              </span>
            </div>

            <div class="flex items-center gap-2 mt-4 pt-3 border-t border-border">
              <RouterLink
                :to="`/profil/${item.id}`"
                class="btn-action text-xs px-3 py-1.5 flex-1 text-center"
              >
                Voir le profil
              </RouterLink>
              <button
                type="button"
                class="w-8 h-8 border border-border flex items-center justify-center hover:bg-surface transition-colors"
                :class="{ 'opacity-40 cursor-not-allowed': !canLike }"
                :title="canLike ? 'Aimer' : 'Réservé aux recruteurs'"
                @click="toggleLike(item.id)"
              >
                <Heart
                  class="w-4 h-4"
                  :class="liked.has(item.id) ? 'fill-action text-action' : 'text-text-muted'"
                />
              </button>
              <button
                type="button"
                class="w-8 h-8 border border-border flex items-center justify-center hover:bg-surface transition-colors"
                :class="{ 'opacity-40 cursor-not-allowed': !canContact }"
                :title="canContact ? 'Contacter' : 'Réservé aux recruteurs'"
                @click="canContact && requireAuth()"
              >
                <MessageSquare class="w-4 h-4 text-text-muted" />
              </button>
            </div>
          </div>
        </article>
      </div>

      <p v-if="!pageItems.length" class="text-center text-text-muted font-marianne py-16">
        Aucun profil à afficher.
      </p>

      <!-- Pagination -->
      <nav
        v-if="totalPages > 1"
        class="flex items-center justify-center gap-1 mt-10 font-marianne text-sm"
        aria-label="Pagination"
      >
        <button
          type="button"
          class="px-3 py-1.5 border border-border bg-white disabled:opacity-40 disabled:cursor-not-allowed hover:bg-surface transition-colors"
          :disabled="page === 1"
          @click="goToPage(page - 1)"
        >
          Précédent
        </button>
        <button
          v-for="n in totalPages"
          :key="n"
          type="button"
          class="w-9 py-1.5 border transition-colors"
          :class="n === page
            ? 'bg-primary text-white border-primary'
            : 'bg-white border-border hover:bg-surface'"
          @click="goToPage(n)"
        >
          {{ n }}
        </button>
        <button
          type="button"
          class="px-3 py-1.5 border border-border bg-white disabled:opacity-40 disabled:cursor-not-allowed hover:bg-surface transition-colors"
          :disabled="page === totalPages"
          @click="goToPage(page + 1)"
        >
          Suivant
        </button>
      </nav>
    </div>

    <!-- Modal choix connexion -->
    <div
      v-if="showAuthModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
      @click.self="showAuthModal = false"
    >
      <div class="bg-white p-8 max-w-sm w-full mx-4 rounded-xl shadow-xl text-center">
        <h2 class="text-xl font-marianne font-bold text-primary mb-2">Vous souhaitez interagir ?</h2>
        <p class="font-spectral text-text-muted text-sm mb-8">Connectez-vous pour continuer. Qui êtes-vous ?</p>
        <div class="flex flex-col gap-4">
          <button class="btn-action w-full" @click="chooseType('candidat')">
            Je suis candidat
          </button>
          <button class="btn-secondary w-full" @click="chooseType('recruteur')">
            Je suis recruteur
          </button>
        </div>
        <button class="mt-6 text-xs font-marianne text-text-muted hover:underline" @click="showAuthModal = false">
          Annuler
        </button>
      </div>
    </div>
  </div>
</template>
