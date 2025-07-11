<template>
  <form @submit.prevent="login" class="space-y-2 max-w-sm mx-auto mt-20">
    <input v-model="username" placeholder="Username" class="w-full p-2 border" />
    <input v-model="password" type="password" placeholder="Password" class="w-full p-2 border" />
    <input v-if="totp" v-model="code" placeholder="TOTP" class="w-full p-2 border" />
    <button class="px-4 py-2 bg-blue-600 text-white">Login</button>
  </form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
const username = ref('')
const password = ref('')
const code = ref('')
const totp = ref(false)

async function login() {
  const res = await fetch('/api/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ Username: username.value, Password: password.value, TOTP: code.value })
  })
  if (res.status === 401) {
    alert('login failed')
  } else if (res.ok) {
    window.location.reload()
  }
}
</script>

<style scoped></style>
