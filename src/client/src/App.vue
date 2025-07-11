<template>
  <Login v-if="!loggedIn" />
  <div v-else>
    <div class="flex space-x-2 p-2 bg-gray-800 text-gray-200">
      <button @click="view='chat'">Chat</button>
      <button @click="view='admin'">Admin</button>
    </div>
    <ChatLayout v-if="view==='chat'" />
    <AdminLayout v-else />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import ChatLayout from './components/ChatLayout.vue'
import Login from './components/Login.vue'
import AdminLayout from './components/AdminLayout.vue'

const view = ref('chat')
const loggedIn = ref(document.cookie.includes('session='))
</script>

<style>
</style>
