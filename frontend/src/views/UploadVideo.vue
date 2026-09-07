<script setup lang="ts">
import { ref } from 'vue'
import { ArrowLeft } from 'lucide-vue-next'

const videoInput = ref<HTMLInputElement | null>(null)
const videoUrl = ref<string | null>(null)
const showSuccess = ref(false)

function openFilePicker() {
    videoInput.value?.click()
}

function onFileSelected(event: Event) {
    const input = event.target as HTMLInputElement
    const file = input.files?.[0]
    if (file) {
        videoUrl.value = URL.createObjectURL(file)
    }
}

function saveVideo() {
    showSuccess.value = true
    setTimeout(() => {
        showSuccess.value = false
    }, 3000)
}

</script>

<template>
    <div class="bg-surface min-h-screen">
        <div class="max-w-4xl mx-auto px-6 py-12 w-full">
            <RouterLink to="/profil/1" class="inline-flex items-center gap-2 text-sm text-primary font-marianne font-medium hover:underline mb-8">
                <ArrowLeft class="w-4 h-4" />
                Retour au profil
            </RouterLink>
            <div class="bg-white border border-border p-8 flex flex-col items-center gap-6">
                <div v-if="showSuccess"
                    class="fixed top-6 left-1/2 -translate-x-1/2 z-50 bg-white border border-border px-6 py-4 fot-marianne font-bold text-success shadow-md">
                    Vidéo enregistrée avec succès !
                </div>
                <h2 class="text-xl font-marianne font-bold text-primary text-center mb-6">Ajouter une vidéo de présentation</h2>
                <input type="file" accept="video/*" ref="videoInput" class="hidden" @change="onFileSelected" />
                <button class="btn-action" @click="openFilePicker">Choisir une vidéo</button>
                <video class="w-full max-w-lg" v-if="videoUrl" :src="videoUrl" controls />
                <button class="btn-action" v-if="videoUrl" @click="saveVideo">Enregistrer</button>
            </div>
        </div>
    </div>
</template>