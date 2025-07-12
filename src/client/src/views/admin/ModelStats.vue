<template>
  <div>
    <h1 class="text-xl mb-2">Model Stats</h1>
    <div v-if="loading">Loading...</div>
    <div v-else class="text-sm space-y-1">
      <div><strong>ID:</strong> {{ stats?.modelId }}</div>
      <div><strong>Downloads:</strong> {{ stats?.downloads }}</div>
      <div><strong>Last Modified:</strong> {{ stats?.lastModified }}</div>
      <div><strong>SHA:</strong> {{ stats?.sha }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRoute } from "vue-router";
import { getModelStats, type ModelDetail } from "../../services/models";

const route = useRoute();
const loading = ref(true);
const stats = ref<ModelDetail | null>(null);

async function load() {
  loading.value = true;
  try {
    stats.value = await getModelStats(route.params.id as string);
  } finally {
    loading.value = false;
  }
}

load();
</script>

<style scoped></style>
