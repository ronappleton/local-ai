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
const users = ref([] as any[])

onMounted(async () => {
  const res = await fetch('/api/users')
  if (res.ok) {
    users.value = await res.json()
  }
})
</script>

<style scoped></style>
