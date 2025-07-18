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
      <button
        @click="refresh"
        class="ml-2 text-sm text-blue-400 hover:underline"
      >
        Refresh
      </button>
    </div>
    <div v-if="loading" class="py-2">Loading...</div>
    <table v-else class="text-sm w-full border-collapse">
      <thead>
        <tr class="text-left border-b border-gray-700">
          <th class="py-1">Model ID</th>
          <th class="py-1">Last Modified</th>
          <th class="py-1">Downloads</th>
          <th class="py-1">Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="m in filteredModels"
          :key="m.modelId"
          class="border-b border-gray-800"
        >
          <td class="py-1">
            <RouterLink
              :to="`/admin/models/${m.modelId}`"
              class="text-blue-400 hover:underline"
            >
              <ModelCard>{{ m.modelId }}</ModelCard>
            </RouterLink>
          </td>
          <td class="py-1">{{ m.lastModified }}</td>
          <td class="py-1">{{ m.downloads }}</td>
          <td class="py-1 space-x-2">
            <button
              @click="startDownload(m.modelId)"
              class="text-xs text-green-400 hover:underline"
            >
              Download
            </button>
            <span v-if="progress[m.modelId] >= 0"
              >{{ progress[m.modelId] }}%</span
            >
            <button
              @click="enable(m.modelId)"
              class="text-xs text-blue-400 hover:underline"
            >
              Enable
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
// View for listing all models. Formerly `ManageModels.vue`.
import { ref, computed } from "vue";
import { RouterLink } from "vue-router";
import {
  getAllModels,
  enableModel,
  downloadModel,
  refreshModels,
  type Model,
} from "../../services/models";
import ModelCard from "../../components/model/ModelCard.vue";

const pipelines = [
  "text-generation",
  "text2text-generation",
  "text-classification",
  "code",
  "conversational",
];

const selected = ref(pipelines[0]);
const models = ref<Model[]>([]);
const loading = ref(false);
const search = ref("");
const progress = ref<Record<string, number>>({});
const filteredModels = computed(() => {
  if (!search.value) return models.value;
  return models.value.filter((m) =>
    m.modelId.toLowerCase().includes(search.value.toLowerCase()),
  );
});

async function load() {
  loading.value = true;
  try {
    console.log("[ModelList] load pipeline", selected.value);
    models.value = await getAllModels(selected.value);
    console.log("[ModelList] loaded", models.value.length, "models");
  } catch (err) {
    console.error("[ModelList] load error", err);
    models.value = [];
  } finally {
    loading.value = false;
  }
}

async function refresh() {
  loading.value = true;
  try {
    await refreshModels(selected.value);
    models.value = await getAllModels(selected.value);
    console.log("[ModelList] refreshed", models.value.length, "models");
  } finally {
    loading.value = false;
  }
}

function select(p: string) {
  selected.value = p;
  console.log("[ModelList] select", p);
  load();
}

async function enable(id: string) {
  try {
    console.log("[ModelList] enable", id);
    await enableModel(id);
    alert("Model enabled");
  } catch {
    alert("Failed to enable model");
  }
}

async function startDownload(id: string) {
  progress.value[id] = 0;
  try {
    await downloadModel(id, (pct) => (progress.value[id] = pct));
    progress.value[id] = 100;
    console.log("[ModelList] download complete", id);
  } catch {
    progress.value[id] = -1;
    alert("Download failed");
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
