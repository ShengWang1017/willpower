<script setup>
import { computed } from 'vue';

const props = defineProps({
  goal: {
    type: Object,
    required: true
  }
});

const categoryConfig = {
  I_WILL: {
    label: 'I WILL',
    icon: '✅',
    card: 'from-emerald-50/95 via-white to-teal-100/90 border-emerald-200/70 shadow-emerald-500/15 text-moss-800',
    labelClass: 'text-moss-500',
    headingClass: 'text-moss-800'
  },
  I_WONT: {
    label: "I WON'T",
    icon: '🙅',
    card: 'from-rose-50/95 via-white to-orange-100/90 border-rose-200/70 shadow-rose-500/15 text-rose-800',
    labelClass: 'text-rose-500',
    headingClass: 'text-rose-800'
  },
  I_WANT: {
    label: 'I WANT',
    icon: '✨',
    card: 'from-indigo-50/95 via-white to-purple-100/90 border-indigo-200/70 shadow-indigo-500/15 text-indigo-800',
    labelClass: 'text-indigo-500',
    headingClass: 'text-indigo-800'
  }
};

const category = computed(() => categoryConfig[props.goal.type] || categoryConfig.I_WILL);

const createdAtLabel = computed(() => {
  if (!props.goal.created_at) {
    return '';
  }
  const date = new Date(props.goal.created_at);
  return Number.isNaN(date.getTime()) ? props.goal.created_at : date.toLocaleDateString();
});

const statusPillClass = computed(() =>
  props.goal.status === 'archived'
    ? 'bg-[#dc2626] text-white'
    : 'bg-[#16a34a] text-white'
);
</script>

<template>
  <div
    :class="[
      'flex flex-col gap-4 rounded-2xl p-6 transition-all duration-300 bg-gradient-to-br border backdrop-blur-sm shadow-md hover:-translate-y-1.5 hover:shadow-2xl hover:shadow-black/15',
      category.card
    ]"
  >
    <div class="space-y-1.5">
      <p class="muted-label" :class="category.labelClass">{{ category.label }}</p>
      <div class="flex items-center gap-2">
        <span class="text-lg">{{ category.icon }}</span>
        <h3 class="text-xl font-semibold leading-tight" :class="category.headingClass">{{ goal.title }}</h3>
      </div>
      <span
        :class="[
          'inline-flex items-center rounded-full px-3 py-0.5 text-xs font-semibold uppercase tracking-wide shadow-sm transition-colors duration-300',
          statusPillClass
        ]"
      >
        {{ goal.status === 'archived' ? 'Archived' : 'Active' }}
      </span>
    </div>
    <div v-if="createdAtLabel" class="text-xs text-midnight-400">
      Created {{ createdAtLabel }}
    </div>
  </div>
</template>
