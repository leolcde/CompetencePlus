<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuth } from '../stores/auth'

const router = useRouter()
const { login } = useAuth()

const formData = reactive({
    email: '',
    password: '',
})

const loading = ref(false)
const error = ref('')

async function submit() {
    error.value = ''
    loading.value = true
    try {
        await login(formData.email, formData.password)
        router.push({ name: 'admin-dashboard' })
    } catch (e) {
        error.value = e instanceof Error ? e.message : 'Impossible de contacter le serveur.'
    } finally {
        loading.value = false
    }
}

</script>

<template>
    <div class="min-h-screen py-12 sm:py-20 px-6 bg-surface flex items-center justify-center">
        <div class="max-w-md mx-auto card p-6 sm:p-10">
            <h1 class="text-3xl font-marianne font-bold text-center mb-2">Administration</h1>
            <p class="text-text-muted font-spectral mb-10">
                Accès réservé aux administrateurs de la plateforme.
            </p>
            <p v-if="error" class="mb-6 alert-error">{{ error }}</p>
            <form class="space-y-6" @submit.prevent="submit">
                <div>
                    <label class="field-label">Adresse e-mail</label>
                    <input v-model="formData.email" type="email" class="field" required />
                </div>
                <div>
                    <label class="field-label">Mot de passe</label>
                    <input v-model="formData.password" type="password" class="field" required />
                </div>

                <button type="submit" class="btn-action w-full mt-4" :disabled="loading">
                    {{ loading ? 'Connexion...' : 'Se connecter' }}
                </button>
            </form>
        </div>
    </div>
</template>