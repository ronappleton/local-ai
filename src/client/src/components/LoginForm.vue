<template>
  <form @submit.prevent="submit" class="space-y-4 max-w-md mx-auto" novalidate>
    <h2 id="login-title" class="text-xl font-semibold text-center">Login</h2>
    <div class="flex items-center bg-gray-700 border border-gray-600 rounded">
      <span class="material-icons px-2 text-gray-400" aria-hidden="true">email</span>
      <label for="login-email" class="sr-only">Email</label>
      <input
          id="login-email"
          v-model="email"
          type="email"
          placeholder="Email"
          class="flex-1 p-2 bg-transparent placeholder-gray-400 text-gray-200"
          required
          autocomplete="login-email"
      />
    </div>
    <div class="flex items-center bg-gray-700 border border-gray-600 rounded">
      <span class="material-icons px-2 text-gray-400" aria-hidden="true">vpn_key</span>
      <label for="login-password" class="sr-only">Password</label>
      <input
        id="login-password"
        v-model="password"
        type="password"
        placeholder="Password"
        class="flex-1 p-2 bg-transparent placeholder-gray-400 text-gray-200"
        required
        autocomplete="current-password"
      />
    </div>
    <ul v-if="errors.length" class="text-red-400 text-sm">
      <li v-for="(e, i) in errors" :key="i">{{ e }}</li>
    </ul>
    <button class="w-full py-2 bg-blue-600 text-white rounded hover:bg-blue-500 disabled:opacity-50">
      Login
    </button>
    <p class="text-sm text-center">
      <a href="#" @click.prevent="$emit('register')" class="underline">Need an account? Register</a>
    </p>
    <p class="text-sm text-center">
      <a href="#" @click.prevent="$emit('reset')" class="underline">Forgot password?</a>
    </p>
  </form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
const email = ref('')
const password = ref('')
const code = ref('')
const errors = ref<string[]>([])

const emit = defineEmits<{
  (e: 'success'): void
  (e: 'register'): void
  (e: 'reset'): void
}>()

async function submit() {
  errors.value = []
  if (!email.value) errors.value.push('Email is required')
  if (!password.value) errors.value.push('Password is required')
  if (errors.value.length) return

  const res = await fetch('/api/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      Email: email.value,
      Password: password.value,
    }),
  })
  if (res.status === 401) {
    errors.value.push('Login failed')
  } else if (res.ok) {
    emit('success')
    window.location.reload()
  }
}
</script>

<style scoped></style>
