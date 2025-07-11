<template>
  <div class="flex h-full bg-gray-900 text-gray-200">
    <nav class="bg-gray-800 w-48 p-4 space-y-2">
      <h2 class="text-lg mb-2">Admin</h2>
      <ul class="space-y-1 text-sm">
        <li>
          <button
            @click="page = 'users'"
            :class="linkClass('users')"
            class="w-full text-left"
          >
            Users
          </button>
        </li>
        <li>
          <button
            @click="page = 'models'"
            :class="linkClass('models')"
            class="w-full text-left"
          >
            Manage Models
          </button>
        </li>
      </ul>
    </nav>
    <main class="flex-1 p-4 overflow-y-auto">
      <div v-if="page === 'users'">
        <h1 class="text-xl mb-2">Users</h1>
        <ul>
          <li v-for="u in users" :key="u.id">
            {{ u.username }} <span v-if="!u.verified">(unverified)</span>
          </li>
        </ul>
      </div>
      <ManageModels v-else />
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { fetchUsers, type User } from "../api";
import ManageModels from "./ManageModels.vue";

const page = ref<"users" | "models">("users");
const users = ref<User[]>([]);

function linkClass(p: "users" | "models") {
  return page.value === p ? "font-semibold" : "hover:underline";
}

async function loadUsers() {
  try {
    users.value = await fetchUsers();
  } catch {
    users.value = [];
  }
}

onMounted(loadUsers);
</script>

<style scoped></style>
