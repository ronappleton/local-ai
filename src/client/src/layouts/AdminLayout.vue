<template>
  <div class="flex h-full bg-gray-900 text-gray-200">
    <nav class="bg-gray-800 w-48 p-4 space-y-2">
      <h2 class="text-lg mb-2">Admin</h2>
      <ul class="space-y-1 text-sm">
        <li>
          <RouterLink to="/admin/dashboard" :class="linkClass('/admin/dashboard')" class="block">Dashboard</RouterLink>
        </li>
        <li>
          <RouterLink to="/admin/models" :class="linkClass('/admin/models')" class="block">Manage Models</RouterLink>
        </li>
      </ul>
    </nav>
    <main class="flex-1 p-4 overflow-y-auto">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
// Layout wrapper for all admin pages. Replaces old AdminLayout.vue logic.
import { useRoute, RouterLink, RouterView } from 'vue-router'

const route = useRoute()
function linkClass(path: string) {
  return route.path.startsWith(path) ? 'font-semibold' : 'hover:underline'
}
</script>

<style scoped></style>
