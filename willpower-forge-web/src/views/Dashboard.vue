<script setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import api from '../services/api';
import GoalCard from '../components/GoalCard.vue';
import GoalSummaryChart from '../components/GoalSummaryChart.vue';
import ParticleBackground from '../components/ParticleBackground.vue';
import FluidGradientBackground from '../components/FluidGradientBackground.vue';
import GlowCard from '../components/GlowCard.vue';
import { useAuthStore } from '../store/auth';
import { useGsapAnimations } from '../composables/useGsapAnimations';

const authStore = useAuthStore();
const router = useRouter();
const { t, locale } = useI18n();
const { animateCards, animateOnScroll, floatingAnimation, cleanup } = useGsapAnimations();

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

const goalTypes = computed(() => [
  { label: t('goalTypes.I_WILL'), value: 'I_WILL' },
  { label: t('goalTypes.I_WONT'), value: 'I_WONT' },
  { label: t('goalTypes.I_WANT'), value: 'I_WANT' }
]);

const goalTypeLabels = computed(() => ({
  I_WILL: {
    title: t('goalTypes.I_WILL'),
    accent: 'border-moss-200 bg-moss-50/60 text-moss-700'
  },
  I_WONT: {
    title: t('goalTypes.I_WONT'),
    accent: 'border-rose-200 bg-rose-50/70 text-rose-700'
  },
  I_WANT: {
    title: t('goalTypes.I_WANT'),
    accent: 'border-indigo-200 bg-indigo-50/70 text-indigo-700'
  }
}));

const summaryOptions = computed(() => [
  { label: t('dashboard.today'), value: 'today' },
  { label: t('dashboard.yesterday'), value: 'yesterday' }
]);

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

const goToRecycleBin = () => {
  router.push('/recycle-bin');
};

const switchLanguage = (lang) => {
  locale.value = lang;
  localStorage.setItem('locale', lang);
};

watch(toast, (value) => {
  if (value) {
    setTimeout(() => {
      toast.value = '';
    }, 2500);
  }
});

onMounted(async () => {
  await fetchGoals();
  await fetchGoalSummaries();

  // Apply GSAP animations after content loads
  setTimeout(() => {
    animateCards('.goal-card', { delay: 0.2 });
    animateOnScroll('.chart-section');
    floatingAnimation('.floating-icon');
  }, 100);
});

onUnmounted(() => {
  cleanup();
});

watch(summaryScope, () => {
  fetchGoalSummaries();
});
</script>

<template>
  <div class="min-h-screen relative">
    <FluidGradientBackground />
    <ParticleBackground />
    <header class="bg-white/85 backdrop-blur border-b border-white/70 shadow-sm relative z-10">
      <div class="mx-auto max-w-5xl px-4 py-4 flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-semibold text-midnight-900">{{ t('dashboard.title') }}</h1>
          <p class="text-sm text-midnight-500">{{ t('dashboard.subtitle') }}</p>
        </div>
        <div class="flex items-center gap-4">
          <div class="flex items-center gap-2">
            <button
              @click="switchLanguage('en')"
              :class="[
                'px-2 py-1 text-xs font-medium rounded transition',
                locale === 'en' ? 'bg-moss-500 text-white' : 'text-slate-600 hover:bg-slate-100'
              ]"
            >
              EN
            </button>
            <button
              @click="switchLanguage('zh')"
              :class="[
                'px-2 py-1 text-xs font-medium rounded transition',
                locale === 'zh' ? 'bg-moss-500 text-white' : 'text-slate-600 hover:bg-slate-100'
              ]"
            >
              中文
            </button>
          </div>
          <button
            @click="goToRecycleBin"
            class="text-sm font-semibold text-slate-600 hover:text-slate-700 flex items-center gap-1"
            :title="t('dashboard.recycleBin')"
          >
            🗑️ {{ t('dashboard.recycleBin') }}
          </button>
          <button @click="logout" class="text-sm font-semibold text-moss-600 hover:text-moss-700">{{ t('common.logout') }}</button>
        </div>
      </div>
    </header>

    <main class="mx-auto max-w-5xl px-4 py-10 space-y-8 relative z-10">
      <section class="space-y-6">
        <div class="flex flex-col gap-2">
          <h2 class="text-lg font-bold text-gray-900">{{ t('dashboard.yourGoals') }}</h2>
          <p class="text-sm text-gray-500">{{ t('dashboard.yourGoalsDesc') }}</p>
        </div>
        <div>
          <div v-if="isLoading" class="text-midnight-500">{{ t('dashboard.loadingGoals') }}</div>
          <div
            v-else-if="errorMessage"
            class="rounded-lg border border-rose-200 bg-rose-50/80 px-4 py-3 text-sm text-rose-700"
          >
            {{ errorMessage }}
          </div>
          <div v-else>
            <p v-if="goals.length === 0" class="mb-4 text-sm text-midnight-500">{{ t('dashboard.noGoals') }}</p>
            <div class="grid gap-6 md:grid-cols-3">
              <section
                v-for="(meta, key) in goalTypeLabels"
                :key="key"
                class="surface-section flex flex-col fade-in-delay"
                style="box-shadow: 0 4px 10px rgba(0, 0, 0, 0.05);"
              >
                <header class="border-b border-gray-100 px-5 py-4">
                  <h3 class="text-base font-bold text-gray-900">{{ meta.title }}</h3>
                </header>
                <div class="flex-1 px-5 py-4 space-y-3">
                  <RouterLink
                    v-for="goal in groupedGoals[key]"
                    :key="goal.id"
                    :to="`/goals/${goal.id}`"
                    class="block transition-all duration-200 hover:-translate-y-0.5 goal-card"
                  >
                    <GoalCard :goal="goal" />
                  </RouterLink>
                  <div
                    v-if="groupedGoals[key].length === 0"
                    class="rounded-xl border border-gray-200 bg-gray-50 px-4 py-6 text-center text-sm text-gray-400"
                  >
                    {{ t('dashboard.nothingHere') }}
                  </div>
                </div>
              </section>
            </div>
          </div>
        </div>
      </section>

      <GlowCard glow-color="#2563eb" :intensity="0.6">
        <section class="surface-section p-8 fade-in" style="box-shadow: 0 4px 10px rgba(0, 0, 0, 0.05);">
          <h2 class="text-base font-bold text-gray-900 mb-5">{{ t('dashboard.createNewGoal') }}</h2>
        <div class="grid gap-4 md:grid-cols-3">
          <div class="md:col-span-1">
            <label for="goalType" class="block text-sm font-medium text-gray-700 mb-2">{{ t('dashboard.goalType') }}</label>
            <select
              id="goalType"
              v-model="goalForm.type"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors"
            >
              <option v-for="option in goalTypes" :key="option.value" :value="option.value">{{ option.label }}</option>
            </select>
          </div>
          <div class="md:col-span-2">
            <label for="goalTitle" class="block text-sm font-medium text-gray-700 mb-2">{{ t('dashboard.title_field') }}</label>
            <div class="flex gap-3">
              <input
                id="goalTitle"
                v-model="goalForm.title"
                type="text"
                class="flex-1 rounded-lg border border-gray-300 bg-white px-3 py-2 text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors"
                :placeholder="t('dashboard.placeholder')"
                maxlength="255"
              />
              <button
                type="button"
                @click="createGoal"
                class="btn-primary"
              >
                {{ t('dashboard.addButton') }}
              </button>
            </div>
          </div>
        </div>
          <div v-if="toast" class="mt-6 rounded-lg border border-green-200 bg-green-50 px-4 py-3 text-sm text-green-700">
            {{ toast }}
          </div>
        </section>
      </GlowCard>

      <section class="surface-section p-8 fade-in chart-section" style="box-shadow: 0 4px 10px rgba(0, 0, 0, 0.05);">
        <header class="flex flex-wrap items-center justify-between gap-4 mb-6">
          <div>
            <h2 class="text-base font-bold text-gray-900">{{ t('dashboard.overallProgress') }}</h2>
            <span class="text-sm text-gray-500">{{ t('dashboard.allGoals') }}</span>
          </div>
          <div class="inline-flex items-center gap-2 rounded-lg border border-gray-200 bg-white p-1 shadow-sm">
            <button
              v-for="option in summaryOptions"
              :key="option.value"
              type="button"
              :class="[
                'px-4 py-2 text-sm font-medium transition-all duration-200',
                summaryScope === option.value
                  ? 'bg-blue-600 text-white shadow-sm'
                  : 'text-gray-600 hover:bg-gray-100 hover:text-gray-900'
              ]"
              style="border-radius: 6px;"
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
          <p v-else class="text-sm text-midnight-500">{{ t('dashboard.noCheckIns') }}</p>
        </div>
      </section>
    </main>
  </div>
</template>
