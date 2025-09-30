<script setup>
import { computed, onMounted, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import api from '../services/api';
import CheckInChart from '../components/CheckInChart.vue';

const route = useRoute();
const router = useRouter();

const goal = ref(null);
const checkIns = ref([]);
const isLoading = ref(false);
const errorMessage = ref('');

const goalId = computed(() => route.params.id);

const typeLabels = {
  I_WILL: 'I WILL',
  I_WONT: "I WON'T",
  I_WANT: 'I WANT'
};

const goalTypeLabel = (type) => typeLabels[type] || type;

const statuses = [
  { label: 'Completed', value: 'completed', color: 'bg-moss-500 hover:bg-moss-600' },
  { label: 'Partial', value: 'partial', color: 'bg-amber-500 hover:bg-amber-600' },
  { label: 'Failed', value: 'failed', color: 'bg-rose-500 hover:bg-rose-600' }
];

const todayISO = () => new Date().toISOString().slice(0, 10);

const todayCheckIn = computed(() => checkIns.value.find((item) => item.date === todayISO()));

const isSubmitting = ref(false);
const submitError = ref('');
const submitMessage = ref('');
const reviewNotes = ref('');
const isUpdatingStatus = ref(false);
const statusError = ref('');

const isGoalActive = computed(() => goal.value?.status === 'active');
const statusButtonLabel = computed(() => (isGoalActive.value ? 'Active' : 'Archived'));

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

onMounted(loadGoalDetail);

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
  <div class="min-h-screen">
    <header class="bg-white/85 backdrop-blur border-b border-white/70 shadow-sm">
      <div class="mx-auto max-w-4xl px-4 py-4 flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-semibold text-midnight-900">Goal Detail</h1>
          <p class="text-sm text-midnight-500">Track your progress history</p>
        </div>
        <button @click="goBack" class="text-sm font-semibold text-moss-600 hover:text-moss-700">Back to Dashboard</button>
      </div>
    </header>

    <main class="mx-auto max-w-4xl px-4 py-10 space-y-10">
      <div v-if="isLoading" class="text-center text-midnight-500">Loading goal information...</div>

      <div v-else>
        <div v-if="errorMessage" class="mb-6 rounded-lg border border-rose-200 bg-rose-50/80 px-4 py-3 text-rose-700">
          {{ errorMessage }}
        </div>

        <section v-if="goal" class="surface-section p-8">
          <p class="muted-label mb-3">{{ goalTypeLabel(goal.type) }}</p>
          <h2 class="text-3xl font-semibold text-midnight-900 mb-4 leading-tight">{{ goal.title }}</h2>
          <div class="flex flex-wrap gap-4 text-sm text-midnight-500">
            <span>Created: {{ formatDate(goal.created_at) }}</span>
            <span>Updated: {{ formatDate(goal.updated_at) }}</span>
          </div>

          <div class="mt-6 space-y-2">
            <span class="muted-label">Goal Status</span>
            <button
              type="button"
              class="inline-flex items-center rounded-lg px-4 py-1.5 text-sm font-semibold transition-colors duration-300 focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-70 disabled:cursor-not-allowed"
              :class="[
                isGoalActive
                  ? 'bg-[#16a34a] text-white focus:ring-[#16a34a33]'
                  : 'bg-[#dc2626] text-white focus:ring-[#dc262633]'
              ]"
              :disabled="isUpdatingStatus"
              @click="toggleGoalStatus"
            >
              {{ statusButtonLabel }}
            </button>
            <p v-if="statusError" class="text-xs text-rose-600">{{ statusError }}</p>
          </div>
        </section>

        <section v-if="goal" class="surface-section p-8">
          <header class="mb-4 flex items-start justify-between gap-4">
            <div>
              <h3 class="text-lg font-semibold text-midnight-900">Today&apos;s Check-in</h3>
            </div>
            <div v-if="todayCheckIn" class="rounded-full border border-midnight-100/70 bg-white/70 px-3 py-1 text-xs font-medium text-midnight-600">
              Current: <span class="capitalize">{{ todayCheckIn.status }}</span>
            </div>
          </header>

          <div class="space-y-4">
            <div class="flex flex-wrap gap-2">
              <button
                v-for="item in statuses"
                :key="item.value"
                type="button"
                :disabled="isSubmitting || !isGoalActive"
                :class="[
                  'flex-1 min-w-[120px] rounded-lg px-3 py-2 text-white text-sm font-medium transition-colors duration-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-white/50 disabled:opacity-50 disabled:cursor-not-allowed',
                  item.color,
                  isSubmitting || !isGoalActive ? '' : 'shadow-md shadow-black/10 hover:shadow-lg hover:-translate-y-0.5'
                ]"
                @click="submitCheckIn(item.value)"
              >
                {{ item.label }}
              </button>
            </div>

            <div>
              <label for="reviewNotes" class="block text-sm font-medium text-midnight-500">Notes (optional)</label>
              <textarea
                id="reviewNotes"
                v-model="reviewNotes"
                class="mt-1 w-full rounded-lg border border-midnight-100/80 bg-white/85 px-3 py-2 text-midnight-800 focus:outline-none focus:ring-2 focus:ring-moss-400/60"
                rows="3"
                placeholder="Reflect on today&apos;s progress"
                :disabled="isSubmitting || !isGoalActive"
              ></textarea>
              <p v-if="todayCheckIn?.review_notes" class="mt-1 text-xs text-midnight-400">
                Last note: {{ todayCheckIn.review_notes }}
              </p>
            </div>

            <div v-if="submitError" class="rounded-lg border border-rose-200 bg-rose-50/70 px-3 py-2 text-sm text-rose-700">{{ submitError }}</div>
            <div v-else-if="submitMessage" class="rounded-lg border border-moss-200 bg-moss-50/80 px-3 py-2 text-sm text-moss-700">{{ submitMessage }}</div>
          </div>
        </section>

        <section v-if="checkIns.length" class="surface-section p-8">
          <header class="flex items-center justify-between">
            <div>
              <h3 class="text-lg font-semibold text-midnight-900">Progress Overview</h3>
              <p class="text-sm text-midnight-500">Daily outcome mix across your check-ins.</p>
            </div>
            <span class="muted-label">Status Counts</span>
          </header>
          <div class="mt-6">
            <CheckInChart :check-ins="checkIns" />
          </div>
        </section>

        <section class="surface-section">
          <header class="border-b border-white/70 px-6 py-4">
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-midnight-900">Check-in History</h3>
              <p v-if="submitMessage && todayCheckIn" class="text-sm text-moss-600">{{ submitMessage }}</p>
            </div>
          </header>

          <div v-if="checkIns.length === 0" class="px-6 py-8 text-center text-midnight-500">
            No check-ins recorded yet.
          </div>

          <div v-else class="overflow-x-auto">
            <table class="min-w-full divide-y divide-midnight-100/70">
              <thead class="bg-white/60 text-left text-xs uppercase tracking-wider text-midnight-400">
                <tr>
                  <th class="px-6 py-3 font-medium">Date</th>
                  <th class="px-6 py-3 font-medium">Status</th>
                  <th class="px-6 py-3 font-medium">Review Notes</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-midnight-100/70 text-sm text-midnight-700">
                <tr v-for="item in checkIns" :key="item.id">
                  <td class="px-6 py-3 whitespace-nowrap">{{ formatDate(item.date) }}</td>
                  <td class="px-6 py-3 whitespace-nowrap capitalize">{{ item.status }}</td>
                  <td class="px-6 py-3">{{ item.review_notes || '—' }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>
      </div>
    </main>
  </div>
</template>
