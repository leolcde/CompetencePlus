<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { login, auth } from '../lib/auth'

const router = useRouter()
const formData = reactive({ email: '', password: '' })
const loading = ref(false)
const error = ref('')

async function submit() {
  error.value = ''
  loading.value = true
  try {
    await login(formData.email, formData.password)
    if (auth.user.value?.role !== 'admin') {
      error.value = 'Accès réservé aux administrateurs.'
      return
    }
    router.push({ name: 'admin-dashboard' })
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Impossible de contacter le serveur.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-surface flex items-center justify-center px-6">
    <div class="max-w-md w-full card p-8">
      <h1 class="text-2xl font-marianne font-bold text-primary mb-1 text-center">Administration</h1>
      <p class="text-text-muted font-spectral text-sm mb-8 text-center">Accès réservé aux administrateurs.</p>

      <p v-if="error" class="mb-6 alert-error">{{ error }}</p>

      <form class="space-y-6" @submit.prevent="submit">
        <div>
          <label class="field-label">Adresse e-mail</label>
          <input v-model="formData.email" type="email" class="field" required />
        </div>
        <div>
          <label class="field-label">
            Mot de passe
          </label>
          <input v-model="formData.password" type="password" class="field" required />
        </div>
        <button type="submit" class="btn-action w-full" :disabled="loading">
          {{ loading ? 'Connexion...' : 'Se connecter' }}
        </button>
      </form>
    </div>
  </div>
</template>
