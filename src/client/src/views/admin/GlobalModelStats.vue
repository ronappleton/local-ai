<template>
  <div>
    <h1 class="text-xl mb-2">Global Model Stats</h1>
    <div v-if="loading">Loading...</div>
    <div v-else class="text-sm">
      Total models available: {{ stats?.total_models }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { getGlobalStats, type GlobalStats } from "../../services/models";

const loading = ref(true);
const stats = ref<GlobalStats | null>(null);

async function load() {
  loading.value = true;
  try {
    stats.value = await getGlobalStats();
  } finally {
    loading.value = false;
  }
}

load();
</script>

<style scoped></style>
