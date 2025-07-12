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
    <div class="my-2">
      <input
        v-model="search"
        type="text"
        placeholder="Search models"
        class="w-full p-1 rounded bg-gray-800 border border-gray-700 text-sm"
      />
    </div>
    <div v-if="loading" class="py-2">Loading...</div>
    <ul v-else class="space-y-1 text-sm">
      <li v-for="m in filteredModels" :key="m.modelId" class="truncate flex items-center space-x-2">
        <ModelCard>{{ m.modelId }}</ModelCard>
        <button @click="enable(m.modelId)" class="text-xs text-blue-400 hover:underline">Enable</button>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
// View for listing all models. Formerly `ManageModels.vue`.
import { ref, computed } from 'vue'
import { getAllModels, enableModel, type Model } from '../../services/models'
import ModelCard from '../../components/model/ModelCard.vue'

const pipelines = [
  "text-generation",
  "text2text-generation",
  "text-classification",
  "code",
  "conversational",
];

const selected = ref(pipelines[0]);
const models = ref<Model[]>([])
const loading = ref(false);
const search = ref('');
const filteredModels = computed(() => {
  if (!search.value) return models.value
  return models.value.filter(m =>
    m.modelId.toLowerCase().includes(search.value.toLowerCase())
  )
});

async function load() {
  loading.value = true
  try {
    models.value = await getAllModels(selected.value)
  } catch {
    models.value = []
  } finally {
    loading.value = false
  }
}

function select(p: string) {
  selected.value = p;
  load();
}

async function enable(id: string) {
  try {
    await enableModel(id)
    alert('Model enabled')
  } catch {
    alert('Failed to enable model')
  }
}

function tabClass(p: string) {
  return selected.value === p
    ? "border-b-2 border-blue-500"
    : "hover:text-gray-300";
}

load();
</script>

<style scoped></style>
