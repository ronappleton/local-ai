<template>
  <form
    @submit.prevent="submit"
    class="space-y-4 max-w-md mx-auto"
    novalidate
  >
    <h2 id="auth-title" class="text-xl font-semibold text-center">
      {{ mode === 'register' ? 'Register' : 'Login' }}
    </h2>
    <div>
      <label for="username" class="sr-only">Username</label>
      <input
        id="username"
        v-model="username"
        placeholder="Username"
        class="w-full p-2 rounded bg-gray-700 border border-gray-600 placeholder-gray-400 text-gray-200"
        required
        autocomplete="username"
      />
    </div>
    <div v-if="mode === 'register'">
      <label for="email" class="sr-only">Email</label>
      <input
        id="email"
        v-model="email"
        type="email"
        placeholder="Email"
        class="w-full p-2 rounded bg-gray-700 border border-gray-600 placeholder-gray-400 text-gray-200"
        required
        autocomplete="email"
      />
    </div>
    <div>
      <label for="password" class="sr-only">Password</label>
      <input
        id="password"
        v-model="password"
        type="password"
        placeholder="Password"
        class="w-full p-2 rounded bg-gray-700 border border-gray-600 placeholder-gray-400 text-gray-200"
        required
        autocomplete="current-password"
      />
    </div>
    <div v-if="totp">
      <label for="totp" class="sr-only">TOTP</label>
      <input
        id="totp"
        v-model="code"
        placeholder="TOTP"
        class="w-full p-2 rounded bg-gray-700 border border-gray-600 placeholder-gray-400 text-gray-200"
      />
    </div>

    <ul v-if="errors.length" class="text-red-400 text-sm">
      <li v-for="(e, i) in errors" :key="i">{{ e }}</li>
    </ul>

    <button
      class="w-full py-2 bg-blue-600 text-white rounded hover:bg-blue-500 disabled:opacity-50"
    >
      {{ mode === 'register' ? 'Register' : 'Login' }}
    </button>
    <p class="text-sm text-center">
      <a href="#" @click.prevent="toggle" class="underline">{{
        mode === 'register'
          ? 'Already have an account? Login'
          : 'Need an account? Register'
      }}</a>
    </p>
  </form>
</template>

<script setup lang="ts">
import { ref } from "vue";
const username = ref("");
const email = ref("");
const password = ref("");
const code = ref("");
const totp = ref(false);
const mode = ref<"login" | "register">("login");
const errors = ref<string[]>([]);
const emit = defineEmits<{
  (e: 'success'): void
}>();

async function submit() {
  errors.value = [];
  if (!username.value) errors.value.push('Username is required');
  if (!password.value) errors.value.push('Password is required');
  if (mode.value === 'register') {
    if (!email.value) {
      errors.value.push('Email is required');
    } else if (!/^\S+@\S+\.\S+$/.test(email.value)) {
      errors.value.push('Email must be valid');
    }
  }
  if (errors.value.length) return;

  if (mode.value === 'login') {
    const res = await fetch('/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        Username: username.value,
        Password: password.value,
        TOTP: code.value,
      }),
    });
    if (res.status === 401) {
      errors.value.push('Login failed');
    } else if (res.ok) {
      emit('success');
      window.location.reload();
    }
  } else {
    const res = await fetch('/api/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        Username: username.value,
        Email: email.value,
        Password: password.value,
      }),
    });
    if (!res.ok) {
      errors.value.push('Register failed');
    } else {
      mode.value = 'login';
    }
  }
}

function toggle() {
  mode.value = mode.value === "login" ? "register" : "login";
}
</script>

<style scoped></style>
