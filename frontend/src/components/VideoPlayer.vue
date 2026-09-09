<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ url: string }>()

const embed = computed(() => {
  const url = props.url?.trim()
  if (!url) return { kind: 'none' as const, src: '' }

  if (/\.(mp4|webm|mov|ogg)$/i.test(url) || url.startsWith('/uploads/')) {
    return { kind: 'file' as const, src: url }
  }

  try {
    const u = new URL(url)
    const host = u.hostname.replace(/^www\./, '')

    if (host === 'youtube.com' || host === 'm.youtube.com') {
      const id = u.searchParams.get('v')
      if (id) return { kind: 'iframe' as const, src: `https://www.youtube.com/embed/${id}` }
    }
    if (host === 'youtu.be') {
      const id = u.pathname.slice(1)
      if (id) return { kind: 'iframe' as const, src: `https://www.youtube.com/embed/${id}` }
    }
    if (host === 'vimeo.com') {
      const id = u.pathname.split('/').filter(Boolean)[0]
      if (id && /^\d+$/.test(id)) return { kind: 'iframe' as const, src: `https://player.vimeo.com/video/${id}` }
    }
  } catch {
    return { kind: 'link' as const, src: url }
  }

  return { kind: 'link' as const, src: url }
})
</script>

<template>
  <div class="w-full aspect-video bg-black rounded-md overflow-hidden">
    <video
      v-if="embed.kind === 'file'"
      :src="embed.src"
      controls
      class="w-full h-full"
    />
    <iframe
      v-else-if="embed.kind === 'iframe'"
      :src="embed.src"
      class="w-full h-full"
      allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
      allowfullscreen
    />
    <a
      v-else-if="embed.kind === 'link'"
      :href="embed.src"
      target="_blank"
      rel="noopener"
      class="w-full h-full flex items-center justify-center text-white font-marianne text-sm underline"
    >
      Voir la vidéo
    </a>
    <div
      v-else
      class="w-full h-full flex items-center justify-center text-white/60 font-marianne text-sm"
    >
      Aucune vidéo
    </div>
  </div>
</template>
