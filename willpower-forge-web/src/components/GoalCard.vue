<script setup>
import { computed } from 'vue';

const props = defineProps({
  goal: {
    type: Object,
    required: true
  }
});

const typeLabels = {
  I_WILL: 'I WILL',
  I_WONT: "I WON'T",
  I_WANT: 'I WANT'
};

const goalTypeLabel = computed(() => typeLabels[props.goal.type] || props.goal.type);

const createdAtLabel = computed(() => {
  if (!props.goal.created_at) {
    return '';
  }
  const date = new Date(props.goal.created_at);
  return Number.isNaN(date.getTime()) ? props.goal.created_at : date.toLocaleDateString();
});
</script>

<template>
  <div class="bg-white rounded-xl shadow p-5 flex flex-col gap-4 transition hover:shadow-md">
    <div>
      <p class="text-xs uppercase tracking-wide text-slate-400">{{ goalTypeLabel }}</p>
      <h3 class="text-xl font-semibold text-slate-800">{{ goal.title }}</h3>
      <span class="inline-flex items-center mt-1 rounded-full bg-slate-100 px-2 py-0.5 text-xs font-medium text-slate-500">
        {{ goal.status }}
      </span>
    </div>
    <div v-if="createdAtLabel" class="text-xs text-slate-400">
      Created {{ createdAtLabel }}
    </div>
  </div>
</template>
