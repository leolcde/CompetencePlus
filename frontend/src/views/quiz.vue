<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ChevronLeft, ChevronRight, ShieldCheck } from 'lucide-vue-next'
import { useAuth } from '../stores/auth'
import { complete } from '../stores/certification'

interface Question {
  id: number
  content: string
  options: string[]
  weight: number
}

const { user } = useAuth()

const profileId = computed(() => Number(user.value?.id ?? 1))
const questions = ref<Question[]>([])
const answers = ref<Record<number, string>>({})
const current = ref(0)
const error = ref('')
const submitting = ref(false)
const result = ref<{ score: number; badge: boolean } | null>(null)

const total = computed(() => questions.value.length)
const question = computed(() => questions.value[current.value])
const answeredCount = computed(() => Object.keys(answers.value).length)
const progress = computed(() => (total.value ? Math.round((answeredCount.value / total.value) * 100) : 0))
const allAnswered = computed(() => total.value > 0 && answeredCount.value === total.value)

onMounted(loadQuestions)

async function loadQuestions() {
  error.value = ''
  try {
    const res = await fetch('/quiz')
    const body = await res.json()
    if (!res.ok) throw new Error(String(res.status))
    questions.value = (body as Record<string, unknown>[]).map((r) => ({
      id: Number(r.id ?? r.ID ?? 0),
      content: String(r.content ?? r.Content ?? ''),
      options: (r.options ?? r.Options ?? ['Oui', 'Non']) as string[],
      weight: Number(r.weight ?? r.Weight ?? 0),
    }))
    current.value = 0
  } catch {
    error.value = 'Impossible de charger le questionnaire.'
  }
}

async function pick(opt: string) {
  if (!question.value) return
  const qid = question.value.id
  answers.value[qid] = opt
  try {
    await fetch('/quiz/answer', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ profile_id: profileId.value, question_id: qid, options: [opt] }),
    })
  } catch {
    /* la réponse locale reste, la synchro réessaiera à la validation */
  }
  if (current.value < total.value - 1) setTimeout(() => current.value++, 150)
}

function prev() {
  if (current.value > 0) current.value--
}
function next() {
  if (current.value < total.value - 1) current.value++
}

async function restart() {
  try {
    await fetch('/quiz/start', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ profile_id: profileId.value }),
    })
  } catch {
    /* ignore */
  }
  answers.value = {}
  result.value = null
  current.value = 0
}

async function validate() {
  submitting.value = true
  error.value = ''
  try {
    const res = await fetch('/quiz/valider', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ profile_id: profileId.value }),
    })
    const b = (await res.json()) as Record<string, unknown>
    result.value = {
      score: Number(b.total_score ?? b.TotalScore ?? 0),
      badge: Boolean(b.badge_earned ?? b.BadgeEarned ?? false),
    }
    if (user.value?.id) complete(user.value.id, result.value.badge, result.value.score)
  } catch {
    error.value = 'La validation a échoué. Réessayez.'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="flex-1 bg-surface py-10 sm:py-16 px-6">
    <div class="max-w-2xl mx-auto">
      <h1 class="text-3xl sm:text-4xl font-marianne font-black text-primary tracking-tight">
        Certification JEB
      </h1>
      <p class="font-spectral text-text-muted mt-2">
        Un court questionnaire de savoir-être professionnel :
        <strong>communication</strong>, <strong>organisation</strong>, <strong>adaptabilité</strong>.
        Répondez spontanément.
      </p>

      <p v-if="error" class="alert-error mt-6">{{ error }}</p>

      <div v-if="result" class="card p-8 mt-8 text-center">
        <ShieldCheck
          class="w-12 h-12 mx-auto mb-3"
          :class="result.badge ? 'text-action' : 'text-text-muted'"
        />
        <p class="font-marianne font-bold text-primary text-lg">
          {{ result.badge ? 'Badge de certification obtenu' : 'Badge non obtenu' }}
        </p>
        <p class="font-spectral text-text-muted text-sm mt-1">
          Score : {{ result.score }}
        </p>
        <div class="flex flex-col sm:flex-row gap-3 justify-center mt-6">
          <RouterLink :to="{ name: 'dashboard' }" class="btn-action">Mon espace</RouterLink>
          <button type="button" class="btn-secondary" @click="restart">Recommencer</button>
        </div>
      </div>

      <template v-else-if="total">
        <div class="mt-8 mb-4">
          <div class="flex justify-between font-marianne text-xs text-text-muted mb-1.5">
            <span>Question {{ current + 1 }} / {{ total }}</span>
            <span>{{ answeredCount }} répondues</span>
          </div>
          <div class="h-1.5 bg-border rounded-full overflow-hidden">
            <div class="h-full bg-action transition-all" :style="{ width: progress + '%' }" />
          </div>
        </div>

        <div class="card p-6 sm:p-10">
          <p class="font-marianne font-bold text-primary text-lg sm:text-xl leading-snug min-h-[3.5rem]">
            {{ question.content }}
          </p>

          <div class="grid grid-cols-2 gap-3 mt-6">
            <button
              v-for="opt in question.options"
              :key="opt"
              type="button"
              class="border rounded-md py-3 font-marianne font-semibold transition-colors"
              :class="answers[question.id] === opt
                ? 'bg-primary text-white border-primary'
                : 'bg-white text-primary border-primary/25 hover:bg-surface'"
              @click="pick(opt)"
            >
              {{ opt }}
            </button>
          </div>

          <div class="flex items-center justify-between mt-8 pt-5 border-t border-border">
            <button
              type="button"
              class="inline-flex items-center gap-1 text-sm font-marianne text-text-muted disabled:opacity-40"
              :disabled="current === 0"
              @click="prev"
            >
              <ChevronLeft class="w-4 h-4" /> Précédent
            </button>
            <button
              type="button"
              class="inline-flex items-center gap-1 text-sm font-marianne text-text-muted disabled:opacity-40"
              :disabled="current === total - 1"
              @click="next"
            >
              Suivant <ChevronRight class="w-4 h-4" />
            </button>
          </div>
        </div>

        <button
          type="button"
          class="btn-action w-full mt-6"
          :disabled="!allAnswered || submitting"
          @click="validate"
        >
          {{ submitting ? 'Validation…' : allAnswered ? 'Valider le questionnaire' : `Répondez aux ${total - answeredCount} questions restantes` }}
        </button>
      </template>

      <p v-else class="font-spectral text-text-muted mt-8">Aucune question disponible.</p>
    </div>
  </div>
</template>
