<script setup>
import { onMounted, ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import api from '../services/api';
import GoalCard from '../components/GoalCard.vue';
import { useAuthStore } from '../store/auth';

const authStore = useAuthStore();
const router = useRouter();

const goals = ref([]);
const isLoading = ref(false);
const errorMessage = ref('');
const toast = ref('');

const goalForm = ref({
  type: 'I_WILL',
  title: ''
});

const goalTypes = [
  { label: 'I WILL', value: 'I_WILL' },
  { label: "I WON'T", value: 'I_WONT' },
  { label: 'I WANT', value: 'I_WANT' }
];

const fetchGoals = async () => {
  isLoading.value = true;
  errorMessage.value = '';
  try {
    const response = await api.get('/goals');
    goals.value = response.data.data || [];
  } catch (error) {
    errorMessage.value = error.response?.data?.message || 'Failed to fetch goals';
  } finally {
    isLoading.value = false;
  }
};

const createGoal = async () => {
  if (!goalForm.value.title.trim()) {
    toast.value = 'Goal title cannot be empty';
    return;
  }

  try {
    const response = await api.post('/goals', goalForm.value);
    goals.value.unshift(response.data.data);
    goalForm.value = { type: 'I_WILL', title: '' };
    toast.value = 'Goal created successfully';
  } catch (error) {
    toast.value = error.response?.data?.message || 'Unable to create goal';
  }
};

const handleCheckInRecorded = (payload) => {
  toast.value = payload.message;
};

const logout = () => {
  authStore.logout();
  router.push('/login');
};

watch(toast, (value) => {
  if (value) {
    setTimeout(() => {
      toast.value = '';
    }, 2500);
  }
});

onMounted(fetchGoals);
</script>

<template>
  <div class="min-h-screen bg-slate-100">
    <header class="bg-white shadow-sm">
      <div class="mx-auto max-w-5xl px-4 py-4 flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-semibold text-slate-800">Willpower Forge</h1>
          <p class="text-sm text-slate-500">Train your willpower every day</p>
        </div>
        <button @click="logout" class="text-sm font-medium text-indigo-600 hover:text-indigo-800">Logout</button>
      </div>
    </header>

    <main class="mx-auto max-w-5xl px-4 py-8">
      <section class="bg-white rounded-xl shadow p-6 mb-8">
        <h2 class="text-lg font-semibold text-slate-700 mb-4">Create a new goal</h2>
        <div class="grid gap-4 md:grid-cols-3">
          <div class="md:col-span-1">
            <label for="goalType" class="block text-sm font-medium text-slate-600">Goal Type</label>
            <select
              id="goalType"
              v-model="goalForm.type"
              class="mt-1 w-full rounded-md border border-slate-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            >
              <option v-for="option in goalTypes" :key="option.value" :value="option.value">{{ option.label }}</option>
            </select>
          </div>
          <div class="md:col-span-2">
            <label for="goalTitle" class="block text-sm font-medium text-slate-600">Title</label>
            <div class="mt-1 flex gap-3">
              <input
                id="goalTitle"
                v-model="goalForm.title"
                type="text"
                class="flex-1 rounded-md border border-slate-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                placeholder="Describe your goal"
                maxlength="255"
              />
              <button
                type="button"
                @click="createGoal"
                class="rounded-md bg-indigo-600 px-4 py-2 text-white font-medium hover:bg-indigo-700"
              >
                Add
              </button>
            </div>
          </div>
        </div>
      </section>

      <section>
        <div v-if="toast" class="mb-4 rounded-md bg-indigo-50 border border-indigo-200 px-4 py-3 text-indigo-700">
          {{ toast }}
        </div>
        <div v-if="errorMessage" class="mb-4 rounded-md bg-red-50 border border-red-200 px-4 py-3 text-red-700">
          {{ errorMessage }}
        </div>
        <div v-if="isLoading" class="text-center text-slate-500">Loading your goals...</div>
        <div v-else>
          <div v-if="goals.length === 0" class="text-center text-slate-500">No goals yet. Create your first one!</div>
          <div v-else class="grid gap-4 md:grid-cols-2">
            <GoalCard
              v-for="goal in goals"
              :key="goal.id"
              :goal="goal"
              @checkin-recorded="handleCheckInRecorded"
            />
          </div>
        </div>
      </section>
    </main>
  </div>
</template>
