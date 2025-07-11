<template>
  <div class="p-4">
    <h1 class="text-xl mb-2">Admin</h1>
    <ul>
      <li v-for="u in users" :key="u.id">{{ u.username }} <span v-if="!u.verified">(unverified)</span></li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { fetchUsers, type User } from '../api'
const users = ref<User[]>([])

onMounted(async () => {
  try {
    users.value = await fetchUsers()
  } catch {
    users.value = []
  }
})
</script>

<style scoped></style>
