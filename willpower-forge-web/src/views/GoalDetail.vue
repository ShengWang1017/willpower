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
  { label: 'Completed', value: 'completed', color: 'bg-emerald-500 hover:bg-emerald-600' },
  { label: 'Partial', value: 'partial', color: 'bg-amber-500 hover:bg-amber-600' },
  { label: 'Failed', value: 'failed', color: 'bg-rose-500 hover:bg-rose-600' }
];

const todayISO = () => new Date().toISOString().slice(0, 10);

const todayCheckIn = computed(() => checkIns.value.find((item) => item.date === todayISO()));

const isSubmitting = ref(false);
const submitError = ref('');
const submitMessage = ref('');
const reviewNotes = ref('');

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
  if (!goal.value) {
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
  <div class="min-h-screen bg-slate-100">
    <header class="bg-white shadow-sm">
      <div class="mx-auto max-w-4xl px-4 py-4 flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-semibold text-slate-800">Goal Detail</h1>
          <p class="text-sm text-slate-500">Track your progress history</p>
        </div>
        <button @click="goBack" class="text-sm font-medium text-indigo-600 hover:text-indigo-800">Back to Dashboard</button>
      </div>
    </header>

    <main class="mx-auto max-w-4xl px-4 py-8">
      <div v-if="isLoading" class="text-center text-slate-500">Loading goal information...</div>

      <div v-else>
        <div v-if="errorMessage" class="mb-6 rounded-md bg-red-50 border border-red-200 px-4 py-3 text-red-700">
          {{ errorMessage }}
        </div>

        <section v-if="goal" class="mb-8 bg-white rounded-xl shadow p-6">
          <p class="text-xs uppercase tracking-wide text-slate-400 mb-1">{{ goalTypeLabel(goal.type) }}</p>
          <h2 class="text-2xl font-semibold text-slate-800 mb-2">{{ goal.title }}</h2>
          <div class="flex flex-wrap gap-4 text-sm text-slate-500">
            <span class="inline-flex items-center rounded-full bg-slate-100 px-3 py-1 text-xs font-medium text-slate-600">
              Status: {{ goal.status }}
            </span>
            <span>Created: {{ formatDate(goal.created_at) }}</span>
            <span>Updated: {{ formatDate(goal.updated_at) }}</span>
          </div>
        </section>

        <section v-if="goal" class="mb-8 bg-white rounded-xl shadow p-6">
          <header class="mb-4 flex items-start justify-between gap-4">
            <div>
              <h3 class="text-lg font-semibold text-slate-700">Today&apos;s Check-in</h3>
              <p class="text-sm text-slate-500">
                Choose a status for today. You can update it at any time before tomorrow.
              </p>
            </div>
            <div v-if="todayCheckIn" class="rounded-full bg-slate-100 px-3 py-1 text-xs font-medium text-slate-600">
              Current: <span class="capitalize">{{ todayCheckIn.status }}</span>
            </div>
          </header>

          <div class="space-y-4">
            <div class="flex flex-wrap gap-2">
              <button
                v-for="item in statuses"
                :key="item.value"
                type="button"
                :disabled="isSubmitting"
                :class="['flex-1 min-w-[120px] rounded-md px-3 py-2 text-white text-sm font-medium transition-colors', item.color, isSubmitting ? 'opacity-70' : '']"
                @click="submitCheckIn(item.value)"
              >
                {{ item.label }}
              </button>
            </div>

            <div>
              <label for="reviewNotes" class="block text-sm font-medium text-slate-600">Notes (optional)</label>
              <textarea
                id="reviewNotes"
                v-model="reviewNotes"
                class="mt-1 w-full rounded-md border border-slate-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                rows="3"
                placeholder="Reflect on today&apos;s progress"
                :disabled="isSubmitting"
              ></textarea>
              <p v-if="todayCheckIn?.review_notes" class="mt-1 text-xs text-slate-500">
                Last note: {{ todayCheckIn.review_notes }}
              </p>
            </div>

            <div v-if="submitError" class="rounded-md bg-red-50 px-3 py-2 text-sm text-red-600">{{ submitError }}</div>
            <div v-else-if="submitMessage" class="rounded-md bg-emerald-50 px-3 py-2 text-sm text-emerald-600">{{ submitMessage }}</div>
          </div>
        </section>

        <section v-if="checkIns.length" class="mb-8 bg-white rounded-xl shadow p-6">
          <header class="flex items-center justify-between">
            <div>
              <h3 class="text-lg font-semibold text-slate-700">Progress Overview</h3>
              <p class="text-sm text-slate-500">Daily outcome mix across your check-ins.</p>
            </div>
            <span class="text-xs font-medium uppercase tracking-wide text-slate-400">Status Counts</span>
          </header>
          <div class="mt-6">
            <CheckInChart :check-ins="checkIns" />
          </div>
        </section>

        <section class="bg-white rounded-xl shadow">
          <header class="border-b border-slate-100 px-6 py-4">
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-slate-700">Check-in History</h3>
              <p v-if="submitMessage && todayCheckIn" class="text-sm text-emerald-600">{{ submitMessage }}</p>
            </div>
          </header>

          <div v-if="checkIns.length === 0" class="px-6 py-8 text-center text-slate-500">
            No check-ins recorded yet.
          </div>

          <div v-else class="overflow-x-auto">
            <table class="min-w-full divide-y divide-slate-200">
              <thead class="bg-slate-50 text-left text-xs uppercase tracking-wider text-slate-500">
                <tr>
                  <th class="px-6 py-3 font-medium">Date</th>
                  <th class="px-6 py-3 font-medium">Status</th>
                  <th class="px-6 py-3 font-medium">Review Notes</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-200 text-sm text-slate-700">
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
