<template>
  <div>
    <h1 class="text-xl mb-2">Model Detail</h1>
    <div v-if="loading">Loading...</div>
    <div v-else class="text-sm space-y-1">
      <div><strong>ID:</strong> {{ model?.modelId }}</div>
      <div><strong>Last Modified:</strong> {{ model?.lastModified }}</div>
      <div><strong>Downloads:</strong> {{ model?.downloads }}</div>
      <div><strong>SHA:</strong> {{ model?.sha }}</div>
      <div>
        <strong>Tags:</strong>
        <span v-for="t in model?.tags" :key="t" class="mr-1">{{ t }}</span>
      </div>
      <div>
        <strong>Files:</strong>
        <ul class="list-disc list-inside">
          <li v-for="f in model?.files" :key="f">{{ f }}</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRoute } from "vue-router";
import { getModelById, type ModelDetail } from "../../services/models";

const route = useRoute();
const loading = ref(true);
const model = ref<ModelDetail | null>(null);

async function load() {
  loading.value = true;
  try {
    model.value = await getModelById(route.params.id as string);
  } catch {
    model.value = null;
  } finally {
    loading.value = false;
  }
}

load();
</script>

<style scoped></style>
