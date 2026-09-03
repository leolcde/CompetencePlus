<script setup lang="ts">
import { ShieldCheck } from 'lucide-vue-next'

withDefaults(defineProps<{ large?: boolean }>(), { large: false })
</script>

<template>
  <div
    :class="[
      'inline-flex items-center gap-2 bg-success text-white font-marianne font-bold uppercase tracking-wide',
      large ? 'px-4 py-2 text-sm' : 'px-3 py-1 text-xs',
    ]"
  >
    <ShieldCheck :class="large ? 'w-5 h-5' : 'w-4 h-4'" />
    Certifié JEB
  </div>
</template>