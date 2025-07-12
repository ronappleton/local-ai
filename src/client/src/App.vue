<template>
  <div class="flex flex-col h-full">
    <header class="flex justify-between items-center p-2 bg-gray-800 text-gray-200">
      <nav class="space-x-2">
        <RouterLink to="/">Chat</RouterLink>
        <RouterLink v-if="loggedIn && isAdmin" to="/admin/dashboard">Admin</RouterLink>
      </nav>
      <div class="relative">
        <button v-if="!loggedIn" @click="showLogin = true">Login</button>
        <div v-else class="cursor-pointer" @click="dropdownOpen = !dropdownOpen">
          {{ username }}
        </div>
        <div v-if="dropdownOpen" class="absolute right-0 bg-gray-700 text-white mt-1 rounded shadow-md py-1 px-2 space-y-1">
          <RouterLink
            v-if="isAdmin"
            to="/admin/dashboard"
            class="block w-full text-left"
            @click="dropdownOpen=false"
          >Admin</RouterLink>
          <button @click="logout" class="block w-full text-left">Logout</button>
        </div>
      </div>
    </header>

    <router-view class="flex-1" />

    <div
      v-if="showLogin"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      role="dialog"
      aria-modal="true"
      @click.self="closeLogin"
    >
      <div class="relative bg-gray-800 text-gray-200 p-6 rounded-lg shadow-lg w-full max-w-md" aria-labelledby="auth-title">
        <button
          class="absolute top-2 right-2 text-gray-400 hover:text-white"
          @click="closeLogin"
          aria-label="Close login modal"
        >
          &times;
        </button>
        <Login @success="closeLogin" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { fetchUsers, type User } from './api'
import Login from './components/Login.vue'
import { RouterLink } from 'vue-router'

const loggedIn = ref(false)
const dropdownOpen = ref(false)
const showLogin = ref(false)
const username = ref('')
const userId = ref<number | null>(null)
const isAdmin = ref(false)

function closeLogin() {
  showLogin.value = false
}

function handleKey(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    closeLogin()
  }
}

watch(showLogin, (val) => {
  if (val) {
    window.addEventListener('keydown', handleKey)
  } else {
    window.removeEventListener('keydown', handleKey)
  }
})

function parseSession() {
  const match = document.cookie.match(/session=(\d+)/)
  if (match) {
    loggedIn.value = true
    userId.value = parseInt(match[1])
  }
}

async function fetchUsername() {
  if (!loggedIn.value || userId.value === null) return
  try {
    const users = await fetchUsers()
    const u = users.find((user) => user.id === userId.value)
    if (u) {
      username.value = u.username
      isAdmin.value = u.admin === true || u.admin === 1
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
