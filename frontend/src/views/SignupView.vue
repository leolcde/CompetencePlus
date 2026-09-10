<script setup lang="ts">
import { register } from '../lib/auth'
import { useRouter } from 'vue-router'
import { reactive, ref } from 'vue'

const router = useRouter()
const formData = reactive({
  name: '',
  email: '',
  password: '',
  birthday: '',
  role: 'candidate',
  city: '',
  sector: '',
  skills: '',
})
const loading = ref(false)
const error = ref('')

async function submit() {
  error.value = ''
  loading.value = true
  try {
    await register({
      name: formData.name,
      email: formData.email,
      password: formData.password,
      birthday: formData.birthday,
      role: formData.role,
      city: formData.city,
      sector: formData.sector,
      skills: formData.skills
        .split(',')
        .map((s) => s.trim())
        .filter(Boolean),
    })
    router.push({ name: formData.role === 'recruiter' ? 'candidates' : 'quiz' })
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Impossible de contacter le serveur.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="flex-1 py-10 sm:py-16 px-6 bg-surface">
    <div class="max-w-xl mx-auto card p-6 sm:p-10">
      <h1 class="text-3xl font-marianne font-bold text-center mb-2">Inscription</h1>
      <p class="text-text-muted font-spectral mb-10 text-center">Rejoignez la plateforme.</p>

      <p v-if="error" class="mb-6 alert-error">{{ error }}</p>

      <form class="space-y-6" @submit.prevent="submit">
        <div>
          <label class="field-label">Nom complet</label>
          <input v-model="formData.name" type="text" class="field" required />
        </div>

        <div>
          <label class="field-label">Adresse e-mail</label>
          <input v-model="formData.email" type="email" class="field" required />
        </div>

        <div>
          <label class="field-label">Mot de passe</label>
          <input v-model="formData.password" type="password" class="field" required minlength="8" />
          <p class="text-xs text-text-muted mt-1 font-spectral">8 caractères minimum.</p>
        </div>

        <div>
          <label class="field-label">Date de naissance</label>
          <input v-model="formData.birthday" type="date" class="field" required />
        </div>

        <div>
          <label class="field-label">Je m'inscris en tant que</label>
          <select v-model="formData.role" class="field">
            <option value="candidate">Candidat</option>
            <option value="recruiter">Recruteur</option>
          </select>
        </div>

        <div>
          <label class="field-label">Ville</label>
          <input v-model="formData.city" type="text" class="field" />
        </div>

        <div>
          <label class="field-label">Secteur</label>
          <input v-model="formData.sector" type="text" class="field" placeholder="Développement, Design, RH…" />
        </div>

        <div>
          <label class="field-label">Compétences</label>
          <input v-model="formData.skills" type="text" class="field" placeholder="Go, Vue.js, SQL" />
          <p class="text-xs text-text-muted mt-1 font-spectral">Séparées par des virgules.</p>
        </div>

        <button type="submit" class="btn-action w-full mt-8" :disabled="loading">
          {{ loading ? 'Envoi...' : 'Valider mon inscription' }}
        </button>
      </form>
    </div>
  </div>
</template>
