<!-- chosiir ma video -->

<script setup lang="ts">
import { ref } from 'vue'

const videoInput = ref<HTMLInputElement | null>(null)
const videoUrl = ref<string | null>(null)

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
    console.log('enregistrement...', videoUrl.value)
}

</script>

<template>
    <div class="max-w-4xl mx-auto px-6 py-12 w-full">
        <div class="bg-white border border-border p-8 flex flex-col items-center gap-6">
            <h2 class="text-xl font-marianne font-bold text-primary text-center mb-6">Ma vidéo de présentation</h2>
            <input type="file" accept="video/*" ref="videoInput" class="hidden" @change="onFileSelected" />
            <button @click="openFilePicker">Choisir une vidéo</button>
            <video v-if="videoUrl" :src="videoUrl" controls />
            <button v-if="videoUrl" @click="saveVideo">Enregistrer</button>
        </div>
    </div>
</template>