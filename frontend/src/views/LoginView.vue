<script setup lang="ts">

import { login } from '../lib/auth'
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const formData = reactive({ email: '', password: '' })
const loading = ref(false)
const error = ref('')

async function submit() {
    error.value = ''
    loading.value = true
    try {
        await login(formData.email, formData.password)
        router.push({ name: 'account' })
    } catch (e) {
        error.value = e instanceof Error ? e.message : "Impossible de contacter le serveur."
    } finally {
        loading.value = false
    }
}

</script>

<template>
    <div>
        <div>
            <p>hello</p>
        </div>

    </div>
</template>