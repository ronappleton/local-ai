<template>
  <div>
    <div class="border-b border-gray-700 mb-4">
      <nav class="flex space-x-4 text-sm">
        <button
          v-for="p in pipelines"
          :key="p"
          @click="select(p)"
          :class="[tabClass(p), 'pb-1']"
        >
          {{ p }}
        </button>
      </nav>
    </div>
    <div v-if="loading" class="py-2">Loading...</div>
    <ul v-else class="space-y-1 text-sm">
      <li v-for="m in models" :key="m.modelId" class="truncate">
        {{ m.modelId }}
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { fetchModels, type ModelInfo } from "../api";

const pipelines = [
  "text-generation",
  "text2text-generation",
  "text-classification",
  "code",
  "conversational",
];

const selected = ref(pipelines[0]);
const models = ref<ModelInfo[]>([]);
const loading = ref(false);

async function load() {
  loading.value = true;
  try {
    models.value = await fetchModels(selected.value);
  } catch {
    models.value = [];
  } finally {
    loading.value = false;
  }
}

function select(p: string) {
  selected.value = p;
  load();
}

function tabClass(p: string) {
  return selected.value === p
    ? "border-b-2 border-blue-500"
    : "hover:text-gray-300";
}

load();
</script>

<style scoped></style>
