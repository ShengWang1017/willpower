<script setup>
import { computed, ref } from 'vue';
import api from '../services/api';

const props = defineProps({
  goal: {
    type: Object,
    required: true
  }
});

const emit = defineEmits(['checkin-recorded']);

const isSubmitting = ref(false);
const errorMessage = ref('');

const statuses = [
  { label: 'Completed', value: 'completed', color: 'bg-emerald-500 hover:bg-emerald-600' },
  { label: 'Partial', value: 'partial', color: 'bg-amber-500 hover:bg-amber-600' },
  { label: 'Failed', value: 'failed', color: 'bg-rose-500 hover:bg-rose-600' }
];

const typeLabels = {
  I_WILL: 'I WILL',
  I_WONT: "I WON'T",
  I_WANT: 'I WANT'
};

const goalTypeLabel = computed(() => typeLabels[props.goal.type] || props.goal.type);

const submitCheckIn = async (status) => {
  isSubmitting.value = true;
  errorMessage.value = '';
  try {
    const response = await api.post('/checkins', {
      goal_id: props.goal.id,
      status,
      review_notes: ''
    });
    emit('checkin-recorded', {
      message: `Check-in recorded as ${status}`,
      data: response.data.data
    });
  } catch (error) {
    errorMessage.value = error.response?.data?.message || 'Failed to record check-in';
  } finally {
    isSubmitting.value = false;
  }
};
</script>

<template>
  <div class="bg-white rounded-xl shadow p-5 flex flex-col gap-4">
    <div>
      <p class="text-xs uppercase tracking-wide text-slate-400">{{ goalTypeLabel }}</p>
      <h3 class="text-xl font-semibold text-slate-800">{{ goal.title }}</h3>
      <span class="inline-flex items-center mt-1 rounded-full bg-slate-100 px-2 py-0.5 text-xs font-medium text-slate-500">
        {{ goal.status }}
      </span>
    </div>

    <div class="flex flex-col gap-2">
      <p class="text-sm text-slate-500">How did you do today?</p>
      <div class="flex flex-wrap gap-2">
        <button
          v-for="item in statuses"
          :key="item.value"
          type="button"
          :class="['flex-1 min-w-[120px] rounded-md px-3 py-2 text-white text-sm font-medium transition-colors', item.color, isSubmitting ? 'opacity-70' : '']"
          :disabled="isSubmitting"
          @click="submitCheckIn(item.value)"
        >
          {{ item.label }}
        </button>
      </div>
      <p v-if="errorMessage" class="text-xs text-red-600">{{ errorMessage }}</p>
    </div>
  </div>
</template>
