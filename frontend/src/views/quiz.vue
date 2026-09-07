<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useAuth } from '../stores/auth'

interface Question {
  id: number
  content: string
  options: string[]
  weight: number
}

const { user } = useAuth()

const profileId = ref<string>(user.value?.id ?? '1')
const questions = ref<Question[]>([])
const answers = ref<Record<number, string>>({})
const result = ref<unknown>(null)
const log = ref<string[]>([])

const current = ref(0)

function push(msg: string) {
  log.value.unshift(`${new Date().toLocaleTimeString()} — ${msg}`)
}

async function loadQuestions() {
  try {
    const res = await fetch('/quiz')
    const body = await res.json()
    if (!res.ok) throw new Error(JSON.stringify(body))
    questions.value = (body as Record<string, unknown>[]).map((r) => ({
      id: Number(r.id ?? r.ID ?? 0),
      content: String(r.content ?? r.Content ?? ''),
      options: (r.options ?? r.Options ?? []) as string[],
      weight: Number(r.weight ?? r.Weight ?? 0),
    }))
    current.value = 0
    push(`${questions.value.length} questions chargées`)
  } catch (e) {
    push(`error GET /quiz : ${e instanceof Error ? e.message : String(e)}`)
  }
}

async function pick(questionId: number, opt: string) {
  answers.value[questionId] = opt
  try {
    const res = await fetch('/quiz/answer', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        profile_id: Number(profileId.value),
        question_id: questionId,
        options: [opt],
      }),
    })
    push(`réponse q${questionId}=${opt} → ${res.status}`)
  } catch (e) {
    push(`error POST /quiz/answer : ${e instanceof Error ? e.message : String(e)}`)
  }
}

function next() {
  if (current.value < questions.value.length - 1) current.value++
}

function prev() {
  if (current.value > 0) current.value--
}

async function restart() {
  try {
    const res = await fetch('/quiz/start', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ profile_id: Number(profileId.value) }),
    })
    answers.value = {}
    result.value = null
    current.value = 0
    push(`reset → ${res.status}`)
  } catch (e) {
    push(`error POST /quiz/start : ${e instanceof Error ? e.message : String(e)}`)
  }
}

async function validate() {
  try {
    const res = await fetch('/quiz/valider', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ profile_id: Number(profileId.value) }),
    })
    result.value = await res.json()
    push(`validation → ${res.status}`)
  } catch (e) {
    push(`error POST /quiz/valider : ${e instanceof Error ? e.message : String(e)}`)
  }
}

onMounted(loadQuestions)
</script>

<template>
  <div>
    <h1>Quiz</h1>

    <p>
      profile_id :
      <input v-model="profileId" />
      <button @click="restart">Recommencer</button>
    </p>

    <div v-if="questions.length">
      <p>Question {{ current + 1 }} / {{ questions.length }}</p>

      <h2>{{ questions[current].content }}</h2>

      <p v-for="opt in questions[current].options" :key="opt">
        <button @click="pick(questions[current].id, opt)">{{ opt }}</button>
        <span v-if="answers[questions[current].id] === opt"> — choisi</span>
      </p>

      <button @click="prev" :disabled="current === 0">Précédent</button>
      <button @click="next" :disabled="current === questions.length - 1">Suivant</button>
      <button @click="validate">Valider le quiz</button>
    </div>

    <p v-else>Aucune question chargée.</p>

    <h2>Résultat</h2>
    <pre>{{ result ? JSON.stringify(result, null, 2) : '(rien)' }}</pre>

    <h2>Log</h2>
    <pre>{{ log.join('\n') }}</pre>
  </div>
</template>
