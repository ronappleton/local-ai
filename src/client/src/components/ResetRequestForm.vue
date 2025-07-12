<template>
  <form @submit.prevent="submit" class="space-y-4 max-w-md mx-auto" novalidate>
    <h2 id="reset-title" class="text-xl font-semibold text-center">Reset Password</h2>
    <div>
      <label for="reset-email" class="sr-only">Email</label>
      <input
        id="reset-email"
        v-model="email"
        type="email"
        placeholder="Email"
        class="w-full p-2 rounded bg-gray-700 border border-gray-600 placeholder-gray-400 text-gray-200"
        required
      />
    </div>
    <p v-if="message" class="text-green-400 text-center">{{ message }}</p>
    <button class="w-full py-2 bg-blue-600 text-white rounded hover:bg-blue-500">
      Send Reset Email
    </button>
    <p class="text-sm text-center">
      <a href="#" @click.prevent="$emit('login')" class="underline">Back to login</a>
    </p>
  </form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
const email = ref('')
const message = ref('')

const emit = defineEmits<{ (e: 'login'): void }>()

async function submit() {
  const res = await fetch('/api/reset/request', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ Email: email.value }),
  })
  if (res.ok) {
    message.value = 'If that email exists, a reset link has been sent.'
  }
}
</script>

<style scoped></style>
