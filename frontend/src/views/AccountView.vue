<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { AlertTriangle, CheckCircle2, MapPin } from 'lucide-vue-next'
import { auth } from '../lib/auth'
import { api } from '../lib/api'

const user = auth.user
const hasConsent = ref(false)
const loading = ref(true)
const error = ref('')

onMounted(async () => {
    try {
        const consent = await api.consent()
        hasConsent.value = consent.active
    } catch (e) {
        error.value = e instanceof Error ? e.message : 'Erreur de chargement'
    } finally {
        loading.value = false
    }
})
</script>

<template>
    <div class="max-w-4xl mx-auto px-6 py-12 w-full">
        <p v-if="loading" class="mb-4 text-sm text-text-muted font-marianne">Chargement de votre profil…</p>
        <p v-if="error" class="mb-4 alert-error">{{ error }}</p>

        <div class="card">
            <div
                class="w-full aspect-video bg-surface relative border-b border-border flex items-center justify-center">
                <div class="flex flex-col items-center justify-center text-text-muted p-6 text-center">
                    <AlertTriangle class="w-12 h-12 mb-4" />
                    <p class="font-marianne font-bold text-lg mb-2">Vidéo non disponible</p>
                    <p class="font-spectral">Aucune présentation pour ce profil.</p>
                </div>
            </div>

            <div class="p-8 md:p-12">
                <div class="mb-12">
                    <h1 class="text-3xl sm:text-4xl font-marianne font-black text-primary mb-2 tracking-tight">
                        {{ user?.name }}
                    </h1>
                    <p class="text-xl font-marianne text-text-main font-medium mb-6">{{ user?.sector }}</p>
                    <div class="flex items-center gap-2 text-text-muted font-marianne text-sm">
                        <MapPin class="w-4 h-4" />
                        {{ user?.city }}
                    </div>
                </div>

                <div class="pt-8 border-t border-border">
                    <h2 class="text-2xl font-marianne font-bold text-primary mb-6">Compétences</h2>
                    <div class="flex flex-wrap gap-2">
                        <span v-for="(skill, idx) in user?.skills" :key="idx"
                            class="px-4 py-2 bg-surface border border-border text-primary font-marianne text-sm font-medium">
                            {{ skill }}
                        </span>
                        <p v-if="!user?.skills?.length" class="text-text-muted font-spectral text-sm italic">
                            Aucune compétence renseignée.
                        </p>
                    </div>
                </div>
            </div>
        </div>

        <div class="mt-12 p-8 border border-border bg-surface">
            <h2 class="text-xl font-marianne font-bold text-primary mb-6">Gestion de ma vidéo (Zone Privée)</h2>

            <div class="card p-6">
                <div class="flex items-start gap-4 mb-6">
                    <input v-model="hasConsent" type="checkbox" id="consent-checkbox"
                        class="mt-1 w-5 h-5 accent-primary" />
                    <div>
                        <label for="consent-checkbox"
                            class="font-marianne font-bold text-primary block mb-2 cursor-pointer">
                            Consentement à la publication de la vidéo
                        </label>
                        <p class="font-spectral text-sm text-text-main leading-relaxed">
                            J'accepte expressément que Compétences+ diffuse ma vidéo de présentation sur la plateforme
                            à destination des recruteurs. Je consens à l'utilisation de mon image et de ma voix dans le
                            cadre exclusif de la mise en relation emploi.
                        </p>
                    </div>
                </div>

                <div
                    class="flex flex-col sm:flex-row sm:items-center justify-between pt-6 border-t border-border gap-4">
                    <div v-if="hasConsent" class="flex items-center gap-2 text-success font-marianne font-bold text-sm">
                        <CheckCircle2 class="w-5 h-5" />
                        Consentement donné
                    </div>
                    <div v-else class="text-text-muted font-marianne text-sm italic">
                        Aucun consentement actif. La vidéo est masquée.
                    </div>

                    <button :disabled="!hasConsent"
                        class="btn-secondary text-sm disabled:opacity-50 disabled:cursor-not-allowed"
                        @click="hasConsent = false">
                        Révoquer mon consentement
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>
