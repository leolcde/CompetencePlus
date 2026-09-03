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

const bgImage = computed(() => {
  if (route.query.type === 'candidat')
    return 'https://images.unsplash.com/photo-1522071820081-009f0129c71c?w=1600&fit=crop&auto=format'
  if (route.query.type === 'recruteur')
    return 'https://images.unsplash.com/photo-1521791136064-7986c2920216?w=1600&fit=crop&auto=format'
  return null
})

const inputClass =
  'w-full border border-border p-3 rounded-none focus:outline-none focus:border-primary font-spectral'
const labelClass = 'block font-marianne font-bold text-primary mb-2 text-sm'

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
  <div class="flex-1 relative py-10 sm:py-16 px-6 bg-surface overflow-hidden">
    <div
      v-if="bgImage"
      class="absolute inset-0 z-0 pointer-events-none opacity-[0.18]"
      :style="`background-image: url('${bgImage}'); background-size: cover; background-position: center;`"
      aria-hidden="true"
    />
    <div class="relative z-10 max-w-md mx-auto bg-white p-6 sm:p-10 border border-border rounded-xl">
      <h1 class="text-3xl font-marianne font-bold text-primary mb-2">{{ title }}</h1>
      <p class="text-text-muted font-spectral mb-10">
        Accédez à votre espace personnel ProfilsActifs.
      </p>

      <p
        v-if="error"
        class="mb-6 border border-action bg-surface text-action font-marianne text-sm p-3"
      >
        {{ error }}
      </p>

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
              class="w-4 h-4 border-border rounded-none text-action focus:ring-action"
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
          :to="{ name: 'signup' }"
          class="text-primary font-marianne font-bold hover:underline underline-offset-4"
        >
          S'inscrire
        </RouterLink>
      </p>
    </div>
  </div>
</template>
