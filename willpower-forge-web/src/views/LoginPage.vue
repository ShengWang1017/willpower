<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../store/auth';

const router = useRouter();
const authStore = useAuthStore();

const form = ref({
  username: '',
  password: ''
});

const isLoading = ref(false);
const errorMessage = ref('');

const handleSubmit = async () => {
  errorMessage.value = '';
  isLoading.value = true;
  try {
    await authStore.login(form.value);
    router.push('/');
  } catch (error) {
    errorMessage.value = error.response?.data?.message || 'Login failed';
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-100 to-slate-200 p-4">
    <div class="w-full max-w-md bg-white shadow-lg rounded-xl p-8">
      <h1 class="text-2xl font-semibold text-center text-slate-800 mb-6">Welcome back</h1>
      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label for="username" class="block text-sm font-medium text-slate-600">Username</label>
          <input
            id="username"
            v-model="form.username"
            type="text"
            class="mt-1 w-full rounded-md border border-slate-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            required
            minlength="3"
            maxlength="50"
          />
        </div>
        <div>
          <label for="password" class="block text-sm font-medium text-slate-600">Password</label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            class="mt-1 w-full rounded-md border border-slate-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            required
            minlength="8"
            maxlength="100"
          />
        </div>
        <p v-if="errorMessage" class="text-sm text-red-600">{{ errorMessage }}</p>
        <button
          type="submit"
          :disabled="isLoading"
          class="w-full flex justify-center items-center rounded-md bg-indigo-600 py-2 text-white font-medium hover:bg-indigo-700 disabled:opacity-70"
        >
          <span v-if="isLoading" class="animate-pulse">Signing in...</span>
          <span v-else>Sign In</span>
        </button>
      </form>
      <p class="mt-4 text-center text-sm text-slate-600">
        New to Willpower Forge?
        <router-link to="/register" class="text-indigo-600 hover:underline">Create an account</router-link>
      </p>
    </div>
  </div>
</template>
