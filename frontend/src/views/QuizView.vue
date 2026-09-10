<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../lib/api'
import type { BadgeResult, Question } from '../type'

const router = useRouter()

const questions = ref<Question[]>([])
const current = ref(0)
const loading = ref(true)
const submitting = ref(false)
const error = ref('')
const result = ref<BadgeResult | null>(null)

const question = computed(() => questions.value[current.value])
const progress = computed(() =>
  questions.value.length ? Math.round((current.value / questions.value.length) * 100) : 0,
)

onMounted(async () => {
  try {
    await api.quizStart()
    questions.value = await api.questions()
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Impossible de charger le questionnaire.'
  } finally {
    loading.value = false
  }
})

async function answer(choice: string) {
  if (submitting.value || !question.value) return
  submitting.value = true
  error.value = ''
  try {
    await api.quizAnswer(question.value.id, choice)
    if (current.value < questions.value.length - 1) {
      current.value++
    } else {
      result.value = await api.quizValidate()
    }
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Une erreur est survenue.'
  } finally {
    submitting.value = false
  }
}

function goToAccount() {
  router.push({ name: 'account' })
}
</script>

<template>
  <section class="max-w-xl mx-auto px-4 py-12">
    <h1>QUIZ</h1>

    <p v-if="loading" class="text-center text-text-muted font-marianne">Chargement…</p>

    <p v-else-if="error" class="alert-error">{{ error }}</p>

    <div v-else-if="result" class="card p-8 text-center">
      <p class="text-5xl mb-4">{{ result.badge ? 'Bravo !!!' : 'Désoler...' }}</p>
      <h1 class="text-2xl font-marianne font-bold mb-2">
        {{ result.badge ? 'Badge obtenu !' : 'Questionnaire terminé' }}
      </h1>
      <p class="text-text-muted font-spectral mb-8">
        Score : {{ result.score }} / {{ questions.length }}
      </p>
      <button class="btn-action w-full" @click="goToAccount">Accéder à mon espace</button>
    </div>

    <div v-else-if="question" class="card p-8">
      <div class="h-1.5 bg-border rounded-full mb-6 overflow-hidden">
        <div class="h-full bg-primary transition-all" :style="{ width: progress + '%' }" />
      </div>

      <p class="text-sm text-text-muted font-marianne mb-2">
        Question {{ current + 1 }} / {{ questions.length }}
      </p>
      <h1 class="text-xl font-marianne font-bold mb-8">{{ question.content }}</h1>

      <div class="grid grid-cols-2 gap-4">
        <button
          v-for="opt in question.options"
          :key="opt"
          class="btn-secondary py-4"
          :disabled="submitting"
          @click="answer(opt)"
        >
          {{ opt }}
        </button>
      </div>
    </div>
  </section>
</template>
