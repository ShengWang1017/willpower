<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { getDeletedGoals, restoreGoal, permanentDeleteGoal } from '../services/api.js';

const router = useRouter();
const deletedGoals = ref([]);
const loading = ref(true);
const error = ref('');

const categoryConfig = {
  I_WILL: {
    label: 'I WILL',
    icon: '✅',
    labelClass: 'text-moss-500',
    headingClass: 'text-moss-800'
  },
  I_WONT: {
    label: "I WON'T",
    icon: '🙅',
    labelClass: 'text-rose-500',
    headingClass: 'text-rose-800'
  },
  I_WANT: {
    label: 'I WANT',
    icon: '✨',
    labelClass: 'text-indigo-500',
    headingClass: 'text-indigo-800'
  }
};

const fetchDeletedGoals = async () => {
  try {
    loading.value = true;
    error.value = '';
    const response = await getDeletedGoals();
    deletedGoals.value = response.data.data || [];
  } catch (err) {
    error.value = err.response?.data?.message || 'Failed to load deleted goals';
    console.error('Error fetching deleted goals:', err);
  } finally {
    loading.value = false;
  }
};

const handleRestore = async (goalId) => {
  if (!confirm('Are you sure you want to restore this goal?')) {
    return;
  }

  try {
    await restoreGoal(goalId);
    await fetchDeletedGoals();
  } catch (err) {
    error.value = err.response?.data?.message || 'Failed to restore goal';
    console.error('Error restoring goal:', err);
  }
};

const handlePermanentDelete = async (goalId) => {
  if (!confirm('Are you sure you want to permanently delete this goal? This action cannot be undone!')) {
    return;
  }

  try {
    await permanentDeleteGoal(goalId);
    await fetchDeletedGoals();
  } catch (err) {
    error.value = err.response?.data?.message || 'Failed to permanently delete goal';
    console.error('Error permanently deleting goal:', err);
  }
};

const formatDate = (dateString) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return Number.isNaN(date.getTime()) ? dateString : date.toLocaleDateString();
};

const getCategoryConfig = (type) => categoryConfig[type] || categoryConfig.I_WILL;

onMounted(() => {
  fetchDeletedGoals();
});
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 via-white to-zinc-100">
    <div class="container mx-auto px-4 py-8">
      <div class="mb-8 flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold text-slate-800">Recycle Bin</h1>
          <p class="mt-2 text-sm text-slate-600">Deleted goals are kept for 30 days before being permanently removed</p>
        </div>
        <button
          @click="router.push('/')"
          class="rounded-lg bg-slate-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-slate-700"
        >
          Back to Dashboard
        </button>
      </div>

      <div v-if="error" class="mb-6 rounded-lg bg-red-50 p-4 text-red-800">
        {{ error }}
      </div>

      <div v-if="loading" class="text-center py-12">
        <div class="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-current border-r-transparent"></div>
        <p class="mt-4 text-slate-600">Loading...</p>
      </div>

      <div v-else-if="deletedGoals.length === 0" class="text-center py-12">
        <div class="text-6xl mb-4">🗑️</div>
        <p class="text-xl text-slate-600">Recycle bin is empty</p>
      </div>

      <div v-else class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <div
          v-for="goal in deletedGoals"
          :key="goal.id"
          class="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm transition hover:shadow-md"
        >
          <div class="mb-4 space-y-2">
            <p class="text-xs font-semibold uppercase tracking-wide" :class="getCategoryConfig(goal.type).labelClass">
              {{ getCategoryConfig(goal.type).label }}
            </p>
            <div class="flex items-center gap-2">
              <span class="text-lg">{{ getCategoryConfig(goal.type).icon }}</span>
              <h3 class="text-lg font-semibold" :class="getCategoryConfig(goal.type).headingClass">
                {{ goal.title }}
              </h3>
            </div>
          </div>

          <div class="mb-4 space-y-1 text-xs text-slate-500">
            <p>Created: {{ formatDate(goal.created_at) }}</p>
            <p>Deleted: {{ formatDate(goal.deleted_at) }}</p>
          </div>

          <div class="flex gap-2">
            <button
              @click="handleRestore(goal.id)"
              class="flex-1 rounded-lg bg-green-600 px-3 py-2 text-sm font-medium text-white transition hover:bg-green-700"
            >
              Restore
            </button>
            <button
              @click="handlePermanentDelete(goal.id)"
              class="flex-1 rounded-lg bg-red-600 px-3 py-2 text-sm font-medium text-white transition hover:bg-red-700"
            >
              Delete Forever
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
