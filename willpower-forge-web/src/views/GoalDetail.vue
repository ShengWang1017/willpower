<script setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import api from '../services/api';
import { deleteGoal } from '../services/api';
import CheckInChart from '../components/CheckInChart.vue';
import DynamicGoalBackground from '../components/DynamicGoalBackground.vue';
import ShaderProgressRing from '../components/ShaderProgressRing.vue';
import { useGsapAnimations } from '../composables/useGsapAnimations';

const route = useRoute();
const router = useRouter();
const { t } = useI18n();
const { animateOnScroll, cleanup } = useGsapAnimations();

const goal = ref(null);
const checkIns = ref([]);
const isLoading = ref(false);
const errorMessage = ref('');
const isEditing = ref(false);
const editForm = ref({
  type: '',
  title: ''
});

const goalId = computed(() => route.params.id);

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

const goalTypeLabel = (type) => goalTypeLabels.value[type]?.title || type;

const statuses = computed(() => [
  { label: t('goalDetail.completed'), value: 'completed', color: 'bg-moss-500 hover:bg-moss-600', icon: '✅' },
  { label: t('goalDetail.partial'), value: 'partial', color: 'bg-amber-500 hover:bg-amber-600', icon: '⚠️' },
  { label: t('goalDetail.failed'), value: 'failed', color: 'bg-rose-500 hover:bg-rose-600', icon: '❌' }
]);

const todayISO = () => new Date().toISOString().slice(0, 10);

const todayCheckIn = computed(() => checkIns.value.find((item) => item.date === todayISO()));

const isSubmitting = ref(false);
const submitError = ref('');
const submitMessage = ref('');
const reviewNotes = ref('');
const isUpdatingStatus = ref(false);
const statusError = ref('');
const historyPageSize = ref(10);
const historyCurrentPage = ref(1);

const isGoalActive = computed(() => goal.value?.status === 'active');
const statusButtonLabel = computed(() => (isGoalActive.value ? t('goalDetail.active') : t('goalDetail.archived')));

// 计算目标完成度
const goalProgress = computed(() => {
  if (checkIns.value.length === 0) return 0;

  const completedCount = checkIns.value.filter((item) => item.status === 'completed').length;
  const partialCount = checkIns.value.filter((item) => item.status === 'partial').length;

  // 完成计为100%，部分完成计为50%
  const totalScore = completedCount * 1 + partialCount * 0.5;
  const progress = (totalScore / checkIns.value.length) * 100;

  return Math.min(100, Math.round(progress));
});

// 计算最近的表现数据用于动态背景
const recentPerformance = computed(() => {
  if (checkIns.value.length === 0) {
    return {
      consecutiveDays: 0,
      recentCompletionRate: 0,
      totalCheckIns: 0
    };
  }

  // 计算连续完成天数（从今天往前数）
  const today = todayISO();
  const sortedCheckIns = [...checkIns.value].sort((a, b) => (a.date > b.date ? -1 : a.date < b.date ? 1 : 0));

  let consecutiveDays = 0;
  let currentDate = new Date(today);

  for (const checkIn of sortedCheckIns) {
    const checkInDate = checkIn.date;
    const expectedDate = currentDate.toISOString().slice(0, 10);

    if (checkInDate === expectedDate && checkIn.status === 'completed') {
      consecutiveDays++;
      currentDate.setDate(currentDate.getDate() - 1);
    } else if (checkInDate === expectedDate && checkIn.status === 'partial') {
      // 部分完成也算半天，但不继续计数
      consecutiveDays += 0.5;
      break;
    } else {
      break;
    }
  }

  // 计算最近7天的完成率
  const sevenDaysAgo = new Date();
  sevenDaysAgo.setDate(sevenDaysAgo.getDate() - 7);
  const sevenDaysAgoISO = sevenDaysAgo.toISOString().slice(0, 10);

  const recentCheckIns = checkIns.value.filter((item) => item.date >= sevenDaysAgoISO);

  let recentCompletionRate = 0;
  if (recentCheckIns.length > 0) {
    const recentCompleted = recentCheckIns.filter((item) => item.status === 'completed').length;
    const recentPartial = recentCheckIns.filter((item) => item.status === 'partial').length;
    const recentScore = recentCompleted * 1 + recentPartial * 0.5;
    recentCompletionRate = recentScore / recentCheckIns.length;
  }

  return {
    consecutiveDays: Math.floor(consecutiveDays),
    recentCompletionRate,
    totalCheckIns: checkIns.value.length
  };
});

// 分页显示的check-in历史
const paginatedCheckIns = computed(() => {
  const start = (historyCurrentPage.value - 1) * historyPageSize.value;
  const end = start + historyPageSize.value;
  return checkIns.value.slice(start, end);
});

const totalPages = computed(() => {
  return Math.ceil(checkIns.value.length / historyPageSize.value);
});

const canGoPrevPage = computed(() => historyCurrentPage.value > 1);
const canGoNextPage = computed(() => historyCurrentPage.value < totalPages.value);

const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    historyCurrentPage.value = page;
  }
};

const nextPage = () => {
  if (canGoNextPage.value) {
    historyCurrentPage.value++;
  }
};

const prevPage = () => {
  if (canGoPrevPage.value) {
    historyCurrentPage.value--;
  }
};

const formatDate = (value) => {
  if (!value) {
    return '-';
  }
  const parsed = new Date(value);
  return Number.isNaN(parsed.getTime()) ? value : parsed.toLocaleDateString();
};

const sortCheckIns = () => {
  checkIns.value = [...checkIns.value].sort((a, b) => (a.date < b.date ? 1 : a.date > b.date ? -1 : 0));
};

const updateGoalStatus = async (status) => {
  if (!goal.value || goal.value.status === status || isUpdatingStatus.value) {
    return;
  }

  isUpdatingStatus.value = true;
  statusError.value = '';

  try {
    const response = await api.patch(`/goals/${goalId.value}/status`, { status });
    goal.value = response.data.data;
  } catch (error) {
    statusError.value = error.response?.data?.message || 'Unable to update status';
  } finally {
    isUpdatingStatus.value = false;
  }
};

const toggleGoalStatus = () => {
  if (!goal.value) {
    return;
  }
  const nextStatus = isGoalActive.value ? 'archived' : 'active';
  updateGoalStatus(nextStatus);
};

const loadGoalDetail = async () => {
  if (!goalId.value) {
    errorMessage.value = 'Invalid goal id';
    goal.value = null;
    checkIns.value = [];
    return;
  }

  isLoading.value = true;
  errorMessage.value = '';
  submitMessage.value = '';
  submitError.value = '';
  try {
    const [goalResponse, checkInsResponse] = await Promise.all([
      api.get(`/goals/${goalId.value}`),
      api.get('/checkins', { params: { goal_id: goalId.value } })
    ]);

    goal.value = goalResponse.data.data;
    checkIns.value = checkInsResponse.data.data || [];
    sortCheckIns();
  } catch (error) {
    errorMessage.value = error.response?.data?.message || 'Unable to load goal details';
    if (error.response?.status === 404) {
      goal.value = null;
      checkIns.value = [];
    }
  } finally {
    isLoading.value = false;
  }
};

const upsertCheckIn = (record) => {
  const existingIndex = checkIns.value.findIndex((item) => item.id === record.id);
  if (existingIndex !== -1) {
    const updated = [...checkIns.value];
    updated[existingIndex] = record;
    checkIns.value = updated;
  } else {
    checkIns.value = [record, ...checkIns.value];
  }
  sortCheckIns();
};

const submitCheckIn = async (status) => {
  if (!goal.value || !isGoalActive.value) {
    return;
  }

  isSubmitting.value = true;
  submitError.value = '';
  submitMessage.value = '';
  try {
    const response = await api.post('/checkins', {
      goal_id: goal.value.id,
      status,
      review_notes: reviewNotes.value
    });

    const newRecord = response.data.data;
    upsertCheckIn(newRecord);
    reviewNotes.value = newRecord.review_notes || '';
    submitMessage.value = `Check-in recorded as ${status}.`;
  } catch (error) {
    submitError.value = error.response?.data?.message || 'Failed to record check-in';
  } finally {
    isSubmitting.value = false;
  }
};

const goBack = () => {
  router.push('/');
};

const handleDelete = async () => {
  if (!confirm(t('goalDetail.deleteConfirm'))) {
    return;
  }

  try {
    await deleteGoal(goalId.value);
    router.push('/');
  } catch (err) {
    console.error('Error deleting goal:', err);
    errorMessage.value = err.response?.data?.message || 'Failed to delete goal';
  }
};

const startEdit = () => {
  if (!goal.value) return;
  editForm.value = {
    type: goal.value.type,
    title: goal.value.title
  };
  isEditing.value = true;
};

const cancelEdit = () => {
  isEditing.value = false;
  editForm.value = { type: '', title: '' };
};

const saveEdit = async () => {
  if (!editForm.value.title.trim()) {
    errorMessage.value = 'Title cannot be empty';
    return;
  }

  try {
    const response = await api.put(`/goals/${goalId.value}`, editForm.value);
    goal.value = response.data.data;
    isEditing.value = false;
    errorMessage.value = '';
  } catch (err) {
    console.error('Error updating goal:', err);
    errorMessage.value = err.response?.data?.message || 'Failed to update goal';
  }
};

const createRipple = (event) => {
  const button = event.currentTarget;
  const ripple = document.createElement('span');
  const rect = button.getBoundingClientRect();
  const size = Math.max(rect.width, rect.height);
  const x = event.clientX - rect.left - size / 2;
  const y = event.clientY - rect.top - size / 2;

  ripple.style.width = ripple.style.height = `${size}px`;
  ripple.style.left = `${x}px`;
  ripple.style.top = `${y}px`;
  ripple.classList.add('ripple');

  button.appendChild(ripple);

  setTimeout(() => {
    ripple.remove();
  }, 600);
};

onMounted(async () => {
  await loadGoalDetail();

  // Apply animations after content loads
  setTimeout(() => {
    animateOnScroll('.animated-section');
  }, 100);
});

onUnmounted(() => {
  cleanup();
});

watch(goalId, (newId, oldId) => {
  if (newId && newId !== oldId) {
    loadGoalDetail();
  }
});

watch(todayCheckIn, (value) => {
  if (value) {
    reviewNotes.value = value.review_notes || '';
  } else {
    reviewNotes.value = '';
    submitMessage.value = '';
  }
});
</script>

<template>
  <div class="min-h-screen relative">
    <DynamicGoalBackground :recent-performance="recentPerformance" />
    <header class="bg-white/85 backdrop-blur border-b border-white/70 shadow-sm relative z-10">
      <div class="mx-auto max-w-4xl px-4 py-4 flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-semibold text-midnight-900">{{ t('goalDetail.title') }}</h1>
          <p class="text-sm text-midnight-500">{{ t('goalDetail.subtitle') }}</p>
        </div>
        <button @click="goBack" class="text-sm font-semibold text-moss-600 hover:text-moss-700">{{ t('goalDetail.backToDashboard') }}</button>
      </div>
    </header>

    <main class="mx-auto max-w-4xl px-4 py-10 space-y-8 relative z-10">
      <div v-if="isLoading" class="text-center text-midnight-500">{{ t('goalDetail.loadingGoal') }}</div>

      <div v-else>
        <div v-if="errorMessage" class="mb-6 rounded-lg border border-rose-200 bg-rose-50/80 px-4 py-3 text-rose-700">
          {{ errorMessage }}
        </div>

        <!-- Goal Info Card - Light shadow, clean white background -->
        <section v-if="goal" class="card-light p-8 fade-in animated-section">
          <div class="flex items-start justify-between gap-6 mb-6">
            <template v-if="!isEditing">
              <div class="flex-1">
                <span :class="['inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold mb-3', goalTypeLabels[goal.type]?.accent]">
                  {{ goalTypeLabel(goal.type) }}
                </span>
                <h2 class="text-lg font-semibold text-midnight-900 mb-3 leading-tight">{{ goal.title }}</h2>
                <div class="flex flex-wrap gap-3 text-sm text-gray-500">
                  <span>{{ t('goalDetail.created') }}: {{ formatDate(goal.created_at) }}</span>
                  <span>{{ t('goalDetail.updated') }}: {{ formatDate(goal.updated_at) }}</span>
                </div>
              </div>
              <div v-if="checkIns.length > 0" class="flex-shrink-0">
                <ShaderProgressRing :progress="goalProgress" :size="140" />
              </div>
            </template>

            <template v-else>
              <div class="flex-1 space-y-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">{{ t('goalDetail.goalType') }}</label>
                  <select
                    v-model="editForm.type"
                    class="w-full border border-gray-300 bg-white px-3 py-2 text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors"
                    style="border-radius: 8px;"
                  >
                    <option v-for="option in goalTypes" :key="option.value" :value="option.value">{{ option.label }}</option>
                  </select>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">{{ t('goalDetail.goalTitle') }}</label>
                  <input
                    v-model="editForm.title"
                    type="text"
                    class="w-full border border-gray-300 bg-white px-3 py-2 text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors"
                    style="border-radius: 8px;"
                    :placeholder="t('goalDetail.titlePlaceholder')"
                    maxlength="255"
                  />
                </div>
              </div>
            </template>

            <div class="flex gap-3 flex-wrap">
              <template v-if="!isEditing">
                <button type="button" class="btn-primary" @click="startEdit">
                  {{ t('goalDetail.editGoal') }}
                </button>
                <button type="button" class="btn-danger" @click="handleDelete">
                  {{ t('goalDetail.deleteGoal') }}
                </button>
              </template>
              <template v-else>
                <button type="button" class="btn-primary" @click="saveEdit">
                  {{ t('goalDetail.save') }}
                </button>
                <button type="button" class="btn-secondary" @click="cancelEdit">
                  {{ t('goalDetail.cancel') }}
                </button>
              </template>
            </div>
          </div>

          <div class="mt-6 pt-6 border-t border-gray-100">
            <div class="flex items-center gap-3">
              <span class="text-sm font-medium text-gray-700">{{ t('goalDetail.statusLabel') }}:</span>
              <span v-if="isGoalActive" class="status-tag-active">{{ t('goalDetail.active') }}</span>
              <span v-else class="status-tag-archived">{{ t('goalDetail.archived') }}</span>
              <button
                type="button"
                class="text-sm text-blue-600 hover:text-blue-700 font-medium ml-2"
                :disabled="isUpdatingStatus"
                @click="toggleGoalStatus"
              >
                {{ isGoalActive ? t('goalDetail.archive') : t('goalDetail.activate') }}
              </button>
            </div>
            <p v-if="statusError" class="text-xs text-rose-600 mt-2">{{ statusError }}</p>
          </div>
        </section>

        <!-- Check-in Card - Medium shadow, interactive section with header background -->
        <section v-if="goal" class="card-medium fade-in">
          <header class="checkin-header px-8 py-4 flex items-start justify-between gap-4 border-b border-gray-100">
            <div class="flex items-center gap-2">
              <span class="text-lg">✅</span>
              <h3 class="text-lg font-semibold text-gray-900">{{ t('goalDetail.todayCheckIn') }}</h3>
            </div>
            <div v-if="todayCheckIn" class="rounded-full border border-midnight-100/70 bg-white/70 px-3 py-1 text-xs font-medium text-midnight-600">
              {{ t('goalDetail.current') }}: <span class="capitalize">{{ todayCheckIn.status }}</span>
            </div>
          </header>

          <div class="p-8 space-y-4">
            <div class="flex gap-3">
              <button
                v-for="item in statuses"
                :key="item.value"
                type="button"
                :disabled="isSubmitting || !isGoalActive"
                :class="['btn-check-in', item.color]"
                @click="(e) => { createRipple(e); submitCheckIn(item.value); }"
              >
                <span>{{ item.icon }}</span>
                <span>{{ item.label }}</span>
              </button>
            </div>

            <div class="pt-4 border-t border-gray-100">
              <label for="reviewNotes" class="block text-sm font-medium text-gray-700 mb-2">{{ t('goalDetail.notesLabel') }}</label>
              <textarea
                id="reviewNotes"
                v-model="reviewNotes"
                class="w-full border border-gray-300 bg-white px-3 py-2 text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors"
                style="border-radius: 8px;"
                rows="3"
                :placeholder="t('goalDetail.notesPlaceholder')"
                :disabled="isSubmitting || !isGoalActive"
              ></textarea>
              <p v-if="todayCheckIn?.review_notes" class="mt-2 text-xs text-gray-500">
                {{ t('goalDetail.lastNote') }}: {{ todayCheckIn.review_notes }}
              </p>
            </div>

            <div v-if="submitError" class="rounded-lg border border-rose-200 bg-rose-50/70 px-3 py-2 text-sm text-rose-700">{{ submitError }}</div>
            <div v-else-if="submitMessage" class="rounded-lg border border-moss-200 bg-moss-50/80 px-3 py-2 text-sm text-moss-700">{{ submitMessage }}</div>
          </div>
        </section>

        <!-- Progress Overview Card - Heavy shadow, gradient background, data result focus -->
        <section v-if="checkIns.length" class="card-heavy fade-in">
          <header class="flex items-center justify-between mb-6">
            <div>
              <h3 class="text-lg font-bold text-blue-700 flex items-center gap-2">
                <span>📊</span>
                {{ t('goalDetail.progressOverview') }}
              </h3>
              <p class="text-sm text-gray-500 mt-1">{{ t('goalDetail.progressDesc') }}</p>
            </div>
            <span class="muted-label">{{ t('goalDetail.statusCounts') }}</span>
          </header>
          <div class="mt-6">
            <CheckInChart :check-ins="checkIns" />
          </div>
        </section>

        <!-- History Card - Heavy shadow, underlined title -->
        <section class="card-heavy fade-in">
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-lg font-semibold text-gray-900 history-title">
              {{ t('goalDetail.checkInHistory') }}
            </h3>
            <span v-if="checkIns.length > 0" class="text-sm text-gray-500">
              {{ t('goalDetail.total') }}: {{ checkIns.length }} {{ t('goalDetail.records') }}
            </span>
          </div>

          <div v-if="checkIns.length === 0" class="py-16 text-center">
            <div class="mb-6 inline-flex items-center justify-center w-24 h-24 rounded-full bg-gray-100">
              <svg class="w-12 h-12 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
              </svg>
            </div>
            <p class="text-gray-500 font-semibold text-lg mb-2">{{ t('goalDetail.noCheckIns') }}</p>
            <p class="text-sm text-gray-400">{{ t('goalDetail.noCheckInsDesc') }}</p>
          </div>

          <div v-else>
            <div class="space-y-4">
              <div
                v-for="(item, index) in paginatedCheckIns"
                :key="item.id"
                class="flex gap-4 relative"
              >
                <!-- Timeline line -->
                <div
                  v-if="index !== paginatedCheckIns.length - 1"
                  class="absolute top-8 bottom-0 border-l border-dashed border-gray-300"
                  style="left: 5px;"
                ></div>

                <!-- Status dot -->
                <div class="relative z-10 flex-shrink-0">
                  <div
                    :class="[
                      'rounded-full border-2 border-white shadow-sm',
                      item.status === 'completed' ? 'bg-green-500' :
                      item.status === 'partial' ? 'bg-yellow-500' :
                      'bg-red-500'
                    ]"
                    style="width: 10px; height: 10px;"
                  ></div>
                </div>

                <!-- Content -->
                <div class="flex-1 pb-6">
                  <div class="flex items-baseline gap-3 mb-1">
                    <span class="text-sm font-semibold text-gray-900">{{ formatDate(item.date) }}</span>
                    <span
                      :class="[
                        'text-xs font-medium px-2 py-0.5 rounded-full capitalize',
                        item.status === 'completed' ? 'bg-green-100 text-green-700' :
                        item.status === 'partial' ? 'bg-yellow-100 text-yellow-700' :
                        'bg-red-100 text-red-700'
                      ]"
                    >
                      {{ item.status }}
                    </span>
                  </div>
                  <p v-if="item.review_notes" class="text-sm text-gray-600 mt-1">{{ item.review_notes }}</p>
                  <p v-else class="text-sm text-gray-400 italic">{{ t('goalDetail.noNotes') }}</p>
                </div>
              </div>
            </div>

            <!-- Pagination -->
            <div v-if="totalPages > 1" class="mt-8 flex items-center justify-center gap-2">
              <button
                type="button"
                @click="prevPage"
                :disabled="!canGoPrevPage"
                :class="[
                  'px-3 py-1.5 rounded-lg text-sm font-medium transition-colors',
                  canGoPrevPage
                    ? 'bg-white border border-gray-300 text-gray-700 hover:bg-gray-50'
                    : 'bg-gray-100 text-gray-400 cursor-not-allowed'
                ]"
              >
                ← {{ t('goalDetail.previous') }}
              </button>

              <div class="flex items-center gap-1">
                <button
                  v-for="page in totalPages"
                  :key="page"
                  type="button"
                  @click="goToPage(page)"
                  :class="[
                    'w-8 h-8 rounded-lg text-sm font-medium transition-colors',
                    page === historyCurrentPage
                      ? 'bg-blue-600 text-white'
                      : 'bg-white border border-gray-300 text-gray-700 hover:bg-gray-50'
                  ]"
                >
                  {{ page }}
                </button>
              </div>

              <button
                type="button"
                @click="nextPage"
                :disabled="!canGoNextPage"
                :class="[
                  'px-3 py-1.5 rounded-lg text-sm font-medium transition-colors',
                  canGoNextPage
                    ? 'bg-white border border-gray-300 text-gray-700 hover:bg-gray-50'
                    : 'bg-gray-100 text-gray-400 cursor-not-allowed'
                ]"
              >
                {{ t('goalDetail.next') }} →
              </button>
            </div>
          </div>
        </section>
      </div>
    </main>
  </div>
</template>
