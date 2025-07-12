<template>
  <form @submit.prevent="submit" class="space-y-4 max-w-md mx-auto" novalidate>
    <h2 id="register-title" class="text-xl font-semibold text-center">Register</h2>
    <div>
      <label for="reg-username" class="sr-only">Username</label>
      <input
        id="reg-username"
        v-model="username"
        placeholder="Username"
        class="w-full p-2 rounded bg-gray-700 border border-gray-600 placeholder-gray-400 text-gray-200"
        required
        autocomplete="username"
      />
    </div>
    <div class="flex items-center bg-gray-700 border border-gray-600 rounded">
      <span class="material-icons px-2 text-gray-400" aria-hidden="true">email</span>
      <label for="reg-email" class="sr-only">Email</label>
      <input
        id="reg-email"
        v-model="email"
        type="email"
        placeholder="Email"
        class="flex-1 p-2 bg-transparent placeholder-gray-400 text-gray-200"
        required
        autocomplete="email"
      />
    </div>
    <div class="flex items-center bg-gray-700 border border-gray-600 rounded">
      <span class="material-icons px-2 text-gray-400" aria-hidden="true">vpn_key</span>
      <label for="reg-password" class="sr-only">Password</label>
      <input
        id="reg-password"
        v-model="password"
        type="password"
        placeholder="Password"
        class="flex-1 p-2 bg-transparent placeholder-gray-400 text-gray-200"
        required
        autocomplete="new-password"
      />
    </div>
    <ul v-if="errors.length" class="text-red-400 text-sm">
      <li v-for="(e, i) in errors" :key="i">{{ e }}</li>
    </ul>
    <button class="w-full py-2 bg-blue-600 text-white rounded hover:bg-blue-500 disabled:opacity-50">
      Register
    </button>
    <p class="text-sm text-center">
      <a href="#" @click.prevent="$emit('login')" class="underline">Already have an account? Login</a>
    </p>
  </form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
const username = ref('')
const email = ref('')
const password = ref('')
const errors = ref<string[]>([])

const emit = defineEmits<{
  (e: 'success'): void
  (e: 'login'): void
}>()

async function submit() {
  errors.value = []
  if (!username.value) errors.value.push('Username is required')
  if (!password.value) errors.value.push('Password is required')
  if (!email.value) {
    errors.value.push('Email is required')
  } else if (!/^\S+@\S+\.\S+$/.test(email.value)) {
    errors.value.push('Email must be valid')
  }
  if (errors.value.length) return

  const res = await fetch('/api/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      Username: username.value,
      Email: email.value,
      Password: password.value,
    }),
  })
  if (!res.ok) {
    errors.value.push('Register failed')
  } else {
    emit('success')
  }
}
</script>

<style scoped></style>
