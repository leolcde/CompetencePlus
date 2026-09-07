<script setup lang="ts">
import { reactive, ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuth } from '../stores/auth'

const router = useRouter()
const route = useRoute()
const { login } = useAuth()

const formData = reactive({
  email: '',
  password: '',
  remember: false,
})

const loading = ref(false)
const error = ref('')

const title = computed(() => {
  if (route.query.type === 'candidat') return 'Connexion candidat'
  if (route.query.type === 'recruteur') return 'Connexion recruteur'
  return 'Connexion'
})

const inputClass = 'field'
const labelClass = 'field-label'

async function submit() {
  error.value = ''
  loading.value = true
  try {
    await login(formData.email, formData.password)

    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : null
    router.push(redirect ?? { name: 'dashboard' })
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Impossible de contacter le serveur.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="flex-1 py-12 sm:py-20 px-6 bg-surface">
    <div class="max-w-md mx-auto card p-6 sm:p-10">
      <h1 class="text-3xl font-marianne font-bold text-primary mb-2">{{ title }}</h1>
      <p class="text-text-muted font-spectral mb-10">
        Accédez à votre espace personnel ProfilsActifs.
      </p>

      <p v-if="error" class="mb-6 alert-error">{{ error }}</p>

      <form class="space-y-6" @submit.prevent="submit">
        <div>
          <label :class="labelClass">Adresse e-mail</label>
          <input v-model="formData.email" type="email" :class="inputClass" required />
        </div>

        <div>
          <label :class="labelClass">Mot de passe</label>
          <input v-model="formData.password" type="password" :class="inputClass" required />
        </div>

        <div class="flex items-center justify-between">
          <label class="flex items-center gap-2 cursor-pointer">
            <input
              v-model="formData.remember"
              type="checkbox"
              class="w-4 h-4 accent-primary"
            />
            <span class="text-sm font-spectral text-text-main">Se souvenir de moi</span>
          </label>
          <a href="#" class="text-sm font-marianne text-primary hover:underline underline-offset-4">
            Mot de passe oublié ?
          </a>
        </div>

        <button type="submit" class="btn-action w-full mt-4" :disabled="loading">
          {{ loading ? 'Connexion...' : 'Se connecter' }}
        </button>
      </form>

      <p class="text-sm font-spectral text-text-muted mt-8 pt-6 border-t border-border text-center">
        Pas encore de compte ?
        <RouterLink
          :to="{ name: 'signup', query: { type: route.query.type } }"
          class="text-primary font-marianne font-bold hover:underline underline-offset-4"
        >
          S'inscrire
        </RouterLink>
      </p>
    </div>
  </div>
</template>
