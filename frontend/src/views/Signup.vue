<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../stores/auth'

const router = useRouter()
const { login } = useAuth()

const formData = reactive({
  firstName: '',
  lastName: '',
  email: '',
  password: '',
  birthDate: '',
  role: 'candidate',
  sector: '',
  location: '',
  skills: '',
  consent: false,
})

const loading = ref(false)
const error = ref('')

const inputClass =
  'w-full border border-border p-3 rounded-none focus:outline-none focus:border-primary font-spectral'
const labelClass = 'block font-marianne font-bold text-primary mb-2 text-sm'

async function submit() {
  error.value = ''

  if (!formData.consent) {
    error.value = 'Vous devez accepter les conditions.'
    return
  }

  loading.value = true
  try {
    const res = await fetch('/auth/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name: `${formData.firstName} ${formData.lastName}`.trim(),
        mail: formData.email,
        password: formData.password,
        role: formData.role,
        date_of_birth: formData.birthDate,
        sector: formData.sector,
        location: formData.location,
        skills: formData.skills
          .split(',')
          .map((s) => s.trim())
          .filter(Boolean),
      }),
    })

    const body = await res.json().catch(() => ({}))

    if (!res.ok) {
      error.value = body.error ?? `Erreur ${res.status}`
      return
    }

    // inscription réussie -> on enchaîne sur un login pour récupérer le token
    await login(formData.email, formData.password)
    router.push({ name: 'dashboard' })
  } catch {
    error.value = 'Impossible de contacter le serveur.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="flex-1 py-10 sm:py-16 px-6 bg-surface">
    <div class="max-w-xl mx-auto bg-white p-6 sm:p-10 border border-border">
      <h1 class="text-3xl font-marianne font-bold text-primary mb-2">Inscription</h1>
      <p class="text-text-muted font-spectral mb-10">
        Rejoignez ProfilsActifs pour accéder à la plateforme.
      </p>

      <p
        v-if="error"
        class="mb-6 border border-action bg-surface text-action font-marianne text-sm p-3"
      >
        {{ error }}
      </p>

      <form class="space-y-6" @submit.prevent="submit">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label :class="labelClass">Prénom</label>
            <input v-model="formData.firstName" type="text" :class="inputClass" required />
          </div>
          <div>
            <label :class="labelClass">Nom</label>
            <input v-model="formData.lastName" type="text" :class="inputClass" required />
          </div>
        </div>

        <div>
          <label :class="labelClass">Adresse e-mail</label>
          <input v-model="formData.email" type="email" :class="inputClass" required />
        </div>

        <div>
          <label :class="labelClass">Mot de passe</label>
          <input
            v-model="formData.password"
            type="password"
            minlength="8"
            maxlength="72"
            :class="inputClass"
            required
          />
          <p class="text-xs text-text-muted mt-1 font-spectral">8 caractères minimum.</p>
        </div>

        <div>
          <label :class="labelClass">Date de naissance *</label>
          <input v-model="formData.birthDate" type="date" :class="inputClass" required />
          <p class="text-xs text-text-muted mt-1 font-spectral">
            Obligatoire pour valider l'inscription.
          </p>
        </div>

        <div>
          <label :class="labelClass">Je m'inscris en tant que</label>
          <select v-model="formData.role" :class="inputClass">
            <option value="candidate">Candidat</option>
            <option value="recruiter">Recruteur</option>
          </select>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label :class="labelClass">Secteur</label>
            <input v-model="formData.sector" type="text" :class="inputClass" />
          </div>
          <div>
            <label :class="labelClass">Localisation</label>
            <input v-model="formData.location" type="text" :class="inputClass" />
          </div>
        </div>

        <div>
          <label :class="labelClass">Compétences</label>
          <input
            v-model="formData.skills"
            type="text"
            :class="inputClass"
            placeholder="Go, SQL, Vue"
          />
          <p class="text-xs text-text-muted mt-1 font-spectral">Séparées par des virgules.</p>
        </div>

        <div class="pt-4 border-t border-border">
          <label class="flex items-start gap-3 cursor-pointer">
            <input
              v-model="formData.consent"
              type="checkbox"
              class="mt-1 w-5 h-5 border-border rounded-none text-action focus:ring-action"
              required
            />
            <span class="text-sm font-spectral text-text-main leading-snug">
              Je consens au traitement de mes données à caractère personnel dans le cadre du service
              ProfilsActifs et j'accepte les Conditions Générales d'Utilisation.
            </span>
          </label>
        </div>

        <button type="submit" class="btn-action w-full mt-8" :disabled="loading">
          {{ loading ? 'Envoi...' : 'Valider mon inscription' }}
        </button>
      </form>
    </div>
  </div>
</template>
