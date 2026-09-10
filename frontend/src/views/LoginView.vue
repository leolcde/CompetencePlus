<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { login, auth } from '../lib/auth'

const router = useRouter()
const route = useRoute()
const formData = reactive({ email: '', password: '' })
const loading = ref(false)
const error = ref('')

async function submit() {
  error.value = ''
  loading.value = true
  try {
    await login(formData.email, formData.password)
    if (auth.user.value?.role === 'admin') {
      router.push({ name: 'admin-dashboard' })
    } else {
      const redirect = route.query.redirect
      router.push(typeof redirect === 'string' ? redirect : { name: 'account' })
    }
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Impossible de contacter le serveur.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="flex-1 py-10 sm:py-16 px-6 bg-surface flex items-center justify-center min-h-screen">
    <div class="max-w-md w-full card p-6 sm:p-10">
      <h1 class="text-3xl font-marianne font-bold text-primary mb-2 text-center">Connexion</h1>
      <p class="text-text-muted font-spectral mb-10 text-center">
        Accédez à votre espace personnel.
      </p>
      <p v-if="error" class="mb-6 alert-error">{{ error }}</p>
      <form class="space-y-6" @submit.prevent="submit">
        <div>
          <label class="field-label">
            Adresse e-mail
          </label>
          <input v-model="formData.email" type="email" class="field" required />
        </div>
        <div>
          <label class="field-label">
            Mot de passe
          </label>
          <input v-model="formData.password" type="password" class="field" required />
        </div>
        <button type="submit" class="btn-action w-full mt-4" :disabled="loading">
          {{ loading ? 'Connexion...' : 'Se connecter' }}
        </button>
      </form>
      <p class="text-sm font-spectral text-text-muted mt-8 pt-6 border-t border-border text-center">
        Pas encore de compte ?
        <RouterLink :to="{ name: 'signup' }" class="text-primary font-marianne font-bold hover:underline">
          S'inscrire
        </RouterLink>
      </p>
    </div>
  </div>
</template>
