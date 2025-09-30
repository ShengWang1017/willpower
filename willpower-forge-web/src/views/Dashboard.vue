<script setup>
import { computed, onMounted, ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import api from '../services/api';
import GoalCard from '../components/GoalCard.vue';
import GoalSummaryChart from '../components/GoalSummaryChart.vue';
import { useAuthStore } from '../store/auth';

const authStore = useAuthStore();
const router = useRouter();

const goals = ref([]);
const goalSummaries = ref([]);
const isLoading = ref(false);
const errorMessage = ref('');
const toast = ref('');
const summaryScope = ref('today');

const goalForm = ref({
  type: 'I_WILL',
  title: ''
});

const goalTypes = [
  { label: 'I WILL', value: 'I_WILL' },
  { label: "I WON'T", value: 'I_WONT' },
  { label: 'I WANT', value: 'I_WANT' }
];

const goalTypeLabels = {
  I_WILL: {
    title: 'I Will',
    accent: 'border-moss-200 bg-moss-50/60 text-moss-700'
  },
  I_WONT: {
    title: "I Won't",
    accent: 'border-rose-200 bg-rose-50/70 text-rose-700'
  },
  I_WANT: {
    title: 'I Want',
    accent: 'border-indigo-200 bg-indigo-50/70 text-indigo-700'
  }
};

const summaryOptions = [
  { label: 'Today', value: 'today' },
  { label: 'Yesterday', value: 'yesterday' }
];

const groupedGoals = computed(() => {
  const groups = {
    I_WILL: [],
    I_WONT: [],
    I_WANT: []
  };

  goals.value.forEach((goal) => {
    const typeKey = (goal.type || '').toUpperCase();
    if (groups[typeKey]) {
      groups[typeKey].push(goal);
    } else {
      groups.I_WILL.push(goal);
    }
  });

  return groups;
});

const fetchGoals = async () => {
  isLoading.value = true;
  errorMessage.value = '';
  try {
    const response = await api.get('/goals');
    const fetchedGoals = Array.isArray(response.data?.data) ? response.data.data : [];
    console.log('[Dashboard] Fetched goals:', fetchedGoals);
    goals.value = fetchedGoals;
  } catch (error) {
    errorMessage.value = error.response?.data?.message || 'Failed to fetch goals';
    console.error('[Dashboard] Failed fetching goals:', error);
  } finally {
    isLoading.value = false;
  }
};

const fetchGoalSummaries = async () => {
  try {
    const params = {};
    if (summaryScope.value === 'today' || summaryScope.value === 'yesterday') {
      const baseDate = new Date();
      if (summaryScope.value === 'yesterday') {
        baseDate.setDate(baseDate.getDate() - 1);
      }
      params.date = baseDate.toISOString().slice(0, 10);
    }
    const response = await api.get('/checkins/summary', { params });
    const summaries = Array.isArray(response.data?.data) ? response.data.data : [];
    console.log('[Dashboard] Fetched summaries:', summaries, 'params:', params);
    goalSummaries.value = summaries;
  } catch (error) {
    // keep silent fallback to avoid blocking dashboard rendering
    console.warn('[Dashboard] Failed to load goal summaries', error);
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
    await fetchGoalSummaries();
  } catch (error) {
    toast.value = error.response?.data?.message || 'Unable to create goal';
  }
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
onMounted(fetchGoalSummaries);

watch(summaryScope, () => {
  fetchGoalSummaries();
});
</script>

<template>
  <div class="min-h-screen">
    <header class="bg-white/85 backdrop-blur border-b border-white/70 shadow-sm">
      <div class="mx-auto max-w-5xl px-4 py-4 flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-semibold text-midnight-900">Willpower Forge</h1>
          <p class="text-sm text-midnight-500">Train your willpower every day</p>
        </div>
        <button @click="logout" class="text-sm font-semibold text-moss-600 hover:text-moss-700">Logout</button>
      </div>
    </header>

    <main class="mx-auto max-w-5xl px-4 py-10 space-y-10">
      <section class="space-y-6">
        <div class="flex flex-col gap-2">
          <h2 class="text-2xl font-semibold text-midnight-900">Your Goals</h2>
          <p class="text-sm text-midnight-500">A quick glance keeps your intentions front and center.</p>
        </div>
        <div>
          <div v-if="isLoading" class="text-midnight-500">Loading your goals...</div>
          <div
            v-else-if="errorMessage"
            class="rounded-lg border border-rose-200 bg-rose-50/80 px-4 py-3 text-sm text-rose-700"
          >
            {{ errorMessage }}
          </div>
          <div v-else>
            <p v-if="goals.length === 0" class="mb-4 text-sm text-midnight-500">No goals yet. Create your first one!</p>
            <div class="grid gap-6 md:grid-cols-3">
              <section
                v-for="(meta, key) in goalTypeLabels"
                :key="key"
                class="surface-section flex flex-col"
              >
                <header class="border-b border-white/70 px-5 py-4">
                  <h3 class="text-lg font-semibold text-midnight-900">{{ meta.title }}</h3>
                </header>
                <div class="flex-1 px-5 py-4 space-y-3">
                  <RouterLink
                    v-for="goal in groupedGoals[key]"
                    :key="goal.id"
                    :to="`/goals/${goal.id}`"
                    class="block transition-all duration-300 hover:-translate-y-1 hover:scale-[1.01]"
                  >
                    <GoalCard :goal="goal" />
                  </RouterLink>
                  <div
                    v-if="groupedGoals[key].length === 0"
                    :class="['rounded-xl border px-4 py-6 text-center text-sm backdrop-blur-sm', meta.accent]"
                  >
                    Nothing here yet. Add a goal to this category.
                  </div>
                </div>
              </section>
            </div>
          </div>
        </div>
      </section>

      <section class="surface-section p-8">
        <h2 class="text-xl font-semibold text-midnight-900 mb-5">Create a new goal</h2>
        <div class="grid gap-4 md:grid-cols-3">
          <div class="md:col-span-1">
            <label for="goalType" class="block text-sm font-medium text-midnight-500">Goal Type</label>
            <select
              id="goalType"
              v-model="goalForm.type"
              class="mt-1 w-full rounded-lg border border-midnight-100/80 bg-white/80 px-3 py-2 text-midnight-700 focus:outline-none focus:ring-2 focus:ring-moss-400/60"
            >
              <option v-for="option in goalTypes" :key="option.value" :value="option.value">{{ option.label }}</option>
            </select>
          </div>
          <div class="md:col-span-2">
            <label for="goalTitle" class="block text-sm font-medium text-midnight-500">Title</label>
            <div class="mt-1 flex gap-3">
              <input
                id="goalTitle"
                v-model="goalForm.title"
                type="text"
                class="flex-1 rounded-lg border border-midnight-100/80 bg-white/90 px-3 py-2 text-midnight-800 focus:outline-none focus:ring-2 focus:ring-moss-400/60"
                placeholder="Describe your goal"
                maxlength="255"
              />
              <button
                type="button"
                @click="createGoal"
                class="inline-flex items-center justify-center rounded-lg bg-moss-500 px-5 py-2 text-sm font-semibold text-white shadow-md shadow-moss-500/30 transition hover:bg-moss-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-moss-400"
              >
                Add
              </button>
            </div>
          </div>
        </div>
        <div v-if="toast" class="mt-6 rounded-lg border border-moss-200 bg-moss-50/70 px-4 py-3 text-sm text-moss-700">
          {{ toast }}
        </div>
      </section>

      <section class="surface-section p-8">
        <header class="flex flex-wrap items-center justify-between gap-4">
          <div>
            <h2 class="text-xl font-semibold text-midnight-900">Overall Progress</h2>
            <span class="muted-label">All Goals</span>
          </div>
          <div class="inline-flex items-center gap-2 rounded-xl border border-midnight-100/70 bg-white/70 p-1 backdrop-blur">
            <button
              v-for="option in summaryOptions"
              :key="option.value"
              type="button"
              :class="[
                'px-4 py-1.5 text-sm font-medium rounded-lg transition-colors',
                summaryScope === option.value
                  ? 'bg-moss-500 text-white shadow shadow-moss-500/30'
                  : 'text-midnight-500 hover:bg-midnight-50'
              ]"
              @click="summaryScope = option.value"
            >
              {{ option.label }}
            </button>
          </div>
        </header>
        <div class="mt-6">
          <template v-if="goalSummaries.length">
            <GoalSummaryChart :summaries="goalSummaries" />
          </template>
          <p v-else class="text-sm text-midnight-500">No check-ins yet for this period.</p>
        </div>
      </section>
    </main>
  </div>
</template>
