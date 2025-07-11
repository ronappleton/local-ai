<template>
  <div class="flex h-full bg-gray-900 text-gray-200 font-sans">
    <nav
      :class="['bg-gray-800 w-64 flex-shrink-0 flex flex-col p-4 space-y-4 overflow-y-auto transition-transform duration-200 ease-in-out', sidebarOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0 md:static']"
      aria-label="Main sidebar"
    >
      <button
        class="md:hidden text-gray-400 hover:text-white mb-4 self-end"
        aria-label="Close sidebar"
        @click="sidebarOpen = false"
      >
        &#x2715;
      </button>
      <ul class="space-y-2">
        <li>
          <a href="#" class="flex items-center space-x-2 hover:bg-gray-700 p-2 rounded">
            <span class="material-icons" aria-hidden="true">add</span>
            <span>New Chat</span>
          </a>
        </li>
        <li>
          <a href="#" class="flex items-center space-x-2 hover:bg-gray-700 p-2 rounded">
            <span class="material-icons" aria-hidden="true">search</span>
            <span>Search chats</span>
          </a>
        </li>
      </ul>
      <div v-if="loggedIn">
        <h2 class="text-sm text-gray-400 uppercase mt-4 mb-2">Projects</h2>
        <ul class="space-y-1 pl-2">
          <li>
            <a href="#" class="flex items-center space-x-2 hover:bg-gray-700 p-2 rounded">
              <span class="material-icons" aria-hidden="true">folder</span>
              <span>Project Alpha</span>
            </a>
            <ul class="pl-6 mt-1 space-y-1">
              <li>
                <a href="#" class="text-gray-400 hover:text-gray-200">Subtask 1</a>
              </li>
              <li>
                <a href="#" class="text-gray-400 hover:text-gray-200">Subtask 2</a>
              </li>
            </ul>
          </li>
          <li>
            <a href="#" class="flex items-center space-x-2 hover:bg-gray-700 p-2 rounded">
              <span class="material-icons" aria-hidden="true">folder</span>
              <span>Project Beta</span>
            </a>
          </li>
        </ul>
      </div>
      <div>
        <h2 class="text-sm text-gray-400 uppercase mt-4 mb-2">Chats</h2>
        <ul class="pl-2 space-y-1 text-gray-400 text-sm">
          <li><a href="#" class="hover:text-gray-200">Chat 1</a></li>
          <li><a href="#" class="hover:text-gray-200">Chat 2</a></li>
          <li><a href="#" class="hover:text-gray-200">Chat 3</a></li>
        </ul>
      </div>
    </nav>

    <div class="flex-1 flex flex-col">
      <header class="p-4 md:hidden">
        <button @click="sidebarOpen = true" class="text-gray-400 hover:text-white" aria-label="Open sidebar">
          &#9776;
        </button>
      </header>
      <main class="flex-grow flex items-center justify-center p-4">
        <div class="text-center w-full max-w-md space-y-3">
          <form class="flex items-center bg-gray-700 rounded-md p-2" @submit.prevent="submit">
            <label for="prompt" class="sr-only">Ask anything</label>
            <span class="material-icons px-2" aria-hidden="true">mic</span>
            <input id="prompt" v-model="prompt" type="text" placeholder="Ask anything" class="flex-grow bg-transparent outline-none px-2" />
            <button type="button" class="px-2" aria-label="Settings">
              <span class="material-icons" aria-hidden="true">settings</span>
            </button>
          </form>
          <p class="text-sm text-gray-400">Ready when you are.</p>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, defineProps } from 'vue'

defineProps<{ loggedIn: boolean }>()

const sidebarOpen = ref(false)
const prompt = ref('')

function submit() {
  // placeholder for future submit action
  console.log('Prompt submitted:', prompt.value)
}
</script>

<style scoped>
</style>
