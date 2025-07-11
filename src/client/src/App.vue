<template>
  <div class="flex flex-col h-full">
    <header class="flex justify-between items-center p-2 bg-gray-800 text-gray-200">
      <nav class="space-x-2">
        <button @click="view='chat'">Chat</button>
        <button v-if="loggedIn && isAdmin" @click="view='admin'">Admin</button>
      </nav>
      <div class="relative">
        <button v-if="!loggedIn" @click="showLogin = true">Login</button>
        <div v-else class="cursor-pointer" @click="dropdownOpen = !dropdownOpen">
          {{ username }}
        </div>
        <div v-if="dropdownOpen" class="absolute right-0 bg-gray-700 text-white mt-1 rounded shadow-md py-1 px-2 space-y-1">
          <button v-if="isAdmin" @click="view='admin'; dropdownOpen=false" class="block w-full text-left">Admin</button>
          <button @click="logout" class="block w-full text-left">Logout</button>
        </div>
      </div>
    </header>

    <ChatLayout v-if="view==='chat'" class="flex-1" />
    <AdminLayout v-else class="flex-1" />

    <div v-if="showLogin" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
      <div class="bg-white p-4 rounded">
        <Login />
        <button class="mt-2 text-sm" @click="showLogin=false">Cancel</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import ChatLayout from './components/ChatLayout.vue'
import Login from './components/Login.vue'
import AdminLayout from './components/AdminLayout.vue'

const view = ref('chat')
const loggedIn = ref(false)
const dropdownOpen = ref(false)
const showLogin = ref(false)
const username = ref('')
const userId = ref<number | null>(null)
const isAdmin = ref(false)

function parseSession() {
  const match = document.cookie.match(/session=(\d+)/)
  if (match) {
    loggedIn.value = true
    userId.value = parseInt(match[1])
    isAdmin.value = userId.value === 1
  }
}

async function fetchUsername() {
  if (!loggedIn.value || userId.value === null) return
  try {
    const res = await fetch('/api/users')
    if (res.ok) {
      const users = await res.json()
      const u = users.find((u: any) => u.id === userId.value)
      if (u) username.value = u.username
    }
  } catch {}
}

function logout() {
  fetch('/api/logout', { method: 'POST' }).finally(() => {
    window.location.reload()
  })
}

onMounted(() => {
  parseSession()
  fetchUsername()
})
</script>

<style>
</style>
