<template>
  <form @submit.prevent="submit" class="space-y-2 max-w-sm mx-auto mt-20">
    <input
      v-model="username"
      placeholder="Username"
      class="w-full p-2 border"
    />
    <input
      v-if="mode === 'register'"
      v-model="email"
      placeholder="Email"
      class="w-full p-2 border"
    />
    <input
      v-model="password"
      type="password"
      placeholder="Password"
      class="w-full p-2 border"
    />
    <input
      v-if="totp"
      v-model="code"
      placeholder="TOTP"
      class="w-full p-2 border"
    />
    <button class="px-4 py-2 bg-blue-600 text-white">
      {{ mode === "register" ? "Register" : "Login" }}
    </button>
    <p class="text-sm text-center">
      <a href="#" @click.prevent="toggle">{{
        mode === "register"
          ? "Already have an account? Login"
          : "Need an account? Register"
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

async function submit() {
  if (mode.value === "login") {
    const res = await fetch("/api/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        Username: username.value,
        Password: password.value,
        TOTP: code.value,
      }),
    });
    if (res.status === 401) {
      alert("login failed");
    } else if (res.ok) {
      window.location.reload();
    }
  } else {
    const res = await fetch("/api/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        Username: username.value,
        Email: email.value,
        Password: password.value,
      }),
    });
    if (!res.ok) {
      alert("register failed");
    } else {
      alert("registered, awaiting verification");
      mode.value = "login";
    }
  }
}

function toggle() {
  mode.value = mode.value === "login" ? "register" : "login";
}
</script>

<style scoped></style>
