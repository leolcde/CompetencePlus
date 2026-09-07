<script setup lang="ts">
import { ref } from 'vue'

withDefaults(defineProps<{ size?: number; withWordmark?: boolean }>(), {
  size: 36,
  withWordmark: false,
})

// public/logo.png, avec repli SVG si le fichier est absent
const logoSrc = '/logo.png'
const useFallback = ref(false)
</script>

<template>
  <span class="inline-flex items-center gap-2.5 shrink-0">
    <img
      v-if="!useFallback"
      :src="logoSrc"
      alt="ProfilsActifs"
      :width="size"
      :height="size"
      class="block object-contain"
      :style="{ width: size + 'px', height: size + 'px' }"
      @error="useFallback = true"
    />
    <svg
      v-else
      :width="size"
      :height="size"
      viewBox="0 0 64 64"
      role="img"
      aria-label="ProfilsActifs"
    >
      <circle cx="32" cy="32" r="29" fill="none" stroke="#71677C" stroke-width="5" stroke-dasharray="60 122" stroke-dashoffset="-2" />
      <circle cx="32" cy="32" r="29" fill="none" stroke="#C08552" stroke-width="5" stroke-dasharray="60 122" stroke-dashoffset="-64" />
      <path d="M6 40a26 26 0 0 0 52 0" fill="none" stroke="#4A001F" stroke-width="5" stroke-linecap="round" />
      <circle cx="20" cy="30" r="6" fill="#71677C" />
      <path d="M11 46c0-5 4-9 9-9s9 4 9 9z" fill="#71677C" />
      <circle cx="44" cy="30" r="6" fill="#C08552" />
      <path d="M35 46c0-5 4-9 9-9s9 4 9 9z" fill="#C08552" />
      <circle cx="32" cy="26" r="8" fill="#4A001F" />
      <path d="M20 45c0-7 5-12 12-12s12 5 12 12z" fill="#4A001F" />
      <path d="M30 34h4l-1 9h-2z" fill="#fff" />
    </svg>
    <span
      v-if="withWordmark"
      class="font-marianne font-black text-primary tracking-tight"
      :style="{ fontSize: Math.round(size * 0.52) + 'px' }"
    >
      ProfilsActifs
    </span>
  </span>
</template>
