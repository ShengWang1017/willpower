<script setup>
import { computed } from 'vue';
import { Bar } from 'vue-chartjs';
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  BarElement,
  CategoryScale,
  LinearScale
} from 'chart.js';

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale);

const props = defineProps({
  summaries: {
    type: Array,
    default: () => []
  }
});

const palette = {
  completed: {
    border: '#379d64',
    gradient: ['rgba(55, 157, 100, 0.55)', 'rgba(55, 157, 100, 0.18)']
  },
  partial: {
    border: '#f59e0b',
    gradient: ['rgba(245, 158, 11, 0.55)', 'rgba(245, 158, 11, 0.18)']
  },
  failed: {
    border: '#ef4444',
    gradient: ['rgba(239, 68, 68, 0.55)', 'rgba(239, 68, 68, 0.18)']
  }
};

const labels = computed(() => props.summaries.map((item) => item.title || `Goal #${item.goal_id}`));

const datasets = computed(() => {
  const statuses = ['completed', 'partial', 'failed'];
  return statuses.map((status) => ({
    status,
    data: props.summaries.map((item) => item[status] ?? 0)
  }));
});

const chartData = computed(() => ({
  labels: labels.value,
  datasets: datasets.value.map((dataset) => ({
    label: dataset.status.charAt(0).toUpperCase() + dataset.status.slice(1),
    data: dataset.data,
    borderWidth: 2,
    borderRadius: 14,
    borderColor: palette[dataset.status].border,
    backgroundColor: (context) => {
      const { ctx, chartArea } = context.chart;
      if (!chartArea) {
        return palette[dataset.status].gradient[0];
      }
      const gradient = ctx.createLinearGradient(0, chartArea.top, 0, chartArea.bottom);
      const [start, end] = palette[dataset.status].gradient;
      gradient.addColorStop(0, start);
      gradient.addColorStop(1, end);
      return gradient;
    },
    hoverBackgroundColor: palette[dataset.status].border,
    maxBarThickness: 40
  }))
}));

const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    x: {
      stacked: true,
      grid: {
        display: false
      },
      ticks: {
        color: '#475569',
        font: {
          family: 'Inter, system-ui, sans-serif',
          size: 12
        }
      }
    },
    y: {
      stacked: true,
      beginAtZero: true,
      grid: {
        color: 'rgba(148, 163, 184, 0.2)'
      },
      ticks: {
        precision: 0,
        color: '#475569',
        font: {
          family: 'Inter, system-ui, sans-serif'
        }
      }
    }
  },
  plugins: {
    legend: {
      position: 'top',
      labels: {
        usePointStyle: true,
        pointStyle: 'roundedRect',
        padding: 16,
        color: '#1f2937',
        font: {
          family: 'Inter, system-ui, sans-serif',
          size: 12
        }
      }
    },
    tooltip: {
      backgroundColor: '#0f172a',
      titleColor: '#e2e8f0',
      bodyColor: '#cbd5f5',
      cornerRadius: 12,
      padding: 12,
      callbacks: {
        label: (context) => {
          const value = context.parsed.y || 0;
          return `${context.dataset.label}: ${value}`;
        }
      }
    },
    title: {
      display: false
    }
  },
  animation: {
    duration: 900,
    easing: 'easeOutCubic',
    delay: (context) => context.dataIndex * 80
  }
}));
</script>

<template>
  <div class="h-80 w-full">
    <Bar :data="chartData" :options="chartOptions" />
  </div>
</template>
