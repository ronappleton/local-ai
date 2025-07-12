<template>
  <div class="flex flex-col h-full">
    <header
      class="flex justify-between items-center p-2 bg-gray-800 text-gray-200"
    >
      <nav class="space-x-2">
        <RouterLink to="/">Chat</RouterLink>
        <RouterLink to="/admin/dashboard">Admin</RouterLink>
      </nav>
    </header>

    <router-view class="flex-1" />
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from "vue-router";
</script>

<style></style>
