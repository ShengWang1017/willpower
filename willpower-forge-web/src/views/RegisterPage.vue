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
const feedback = ref('');

const handleSubmit = async () => {
  feedback.value = '';
  isLoading.value = true;
  try {
    await authStore.register(form.value);
    feedback.value = 'Registration successful. Redirecting to login...';
    setTimeout(() => router.push('/login'), 1200);
  } catch (error) {
    feedback.value = error.response?.data?.message || 'Registration failed';
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <div class="min-h-screen flex items-center justify-center px-4 py-12">
    <div class="w-full max-w-md surface-section p-8">
      <h1 class="text-3xl font-semibold text-center text-midnight-900 mb-6">Create your account</h1>
      <form @submit.prevent="handleSubmit" class="space-y-5">
        <div class="space-y-1.5">
          <label for="username" class="block text-sm font-medium text-midnight-500">Username</label>
          <input
            id="username"
            v-model="form.username"
            type="text"
            class="w-full rounded-lg border border-midnight-100/80 bg-white/85 px-3 py-2 text-midnight-800 focus:outline-none focus:ring-2 focus:ring-moss-400/60"
            required
            minlength="3"
            maxlength="50"
          />
        </div>
        <div class="space-y-1.5">
          <label for="password" class="block text-sm font-medium text-midnight-500">Password</label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            class="w-full rounded-lg border border-midnight-100/80 bg-white/85 px-3 py-2 text-midnight-800 focus:outline-none focus:ring-2 focus:ring-moss-400/60"
            required
            minlength="8"
            maxlength="100"
          />
        </div>
        <p
          v-if="feedback"
          class="rounded-lg border px-3 py-2 text-sm"
          :class="feedback.includes('successful') ? 'border-moss-200 bg-moss-50/80 text-moss-700' : 'border-rose-200 bg-rose-50/80 text-rose-700'"
        >
          {{ feedback }}
        </p>
        <button
          type="submit"
          :disabled="isLoading"
          class="w-full flex justify-center items-center rounded-lg bg-moss-500 py-2.5 text-white font-semibold shadow-md shadow-moss-500/30 transition hover:bg-moss-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-moss-400 disabled:opacity-60"
        >
          <span v-if="isLoading" class="animate-pulse">Creating...</span>
          <span v-else>Create Account</span>
        </button>
      </form>
      <p class="mt-6 text-center text-sm text-midnight-500">
        Already have an account?
        <router-link to="/login" class="font-semibold text-moss-600 hover:text-moss-700">Sign in here</router-link>
      </p>
    </div>
  </div>
</template>
