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
  checkIns: {
    type: Array,
    required: true
  }
});

const statusPalette = {
  completed: {
    border: '#10b981',
    gradient: ['rgba(16, 185, 129, 0.45)', 'rgba(16, 185, 129, 0.05)']
  },
  partial: {
    border: '#f59e0b',
    gradient: ['rgba(245, 158, 11, 0.45)', 'rgba(245, 158, 11, 0.05)']
  },
  failed: {
    border: '#f87171',
    gradient: ['rgba(248, 113, 113, 0.45)', 'rgba(248, 113, 113, 0.05)']
  }
};

const aggregated = computed(() => {
  const result = new Map();
  [...props.checkIns]
    .sort((a, b) => (a.date < b.date ? -1 : a.date > b.date ? 1 : 0))
    .forEach((item) => {
      const dateKey = item.date;
      if (!result.has(dateKey)) {
        result.set(dateKey, { completed: 0, partial: 0, failed: 0 });
      }
      const group = result.get(dateKey);
      if (group[item.status] !== undefined) {
        group[item.status] += 1;
      }
    });

  const labels = Array.from(result.keys());
  const totals = labels.map((label) => result.get(label));

  return {
    labels,
    datasets: ['completed', 'partial', 'failed'].map((status) => ({
      status,
      data: totals.map((entry) => entry[status])
    }))
  };
});

const chartData = computed(() => ({
  labels: aggregated.value.labels,
  datasets: aggregated.value.datasets.map((dataset) => ({
    label: dataset.status.charAt(0).toUpperCase() + dataset.status.slice(1),
    data: dataset.data,
    borderWidth: 2,
    borderRadius: 12,
    borderColor: statusPalette[dataset.status].border,
    backgroundColor: (context) => {
      const { ctx, chartArea } = context.chart;
      if (!chartArea) {
        return statusPalette[dataset.status].gradient[0];
      }
      const gradient = ctx.createLinearGradient(0, chartArea.top, 0, chartArea.bottom);
      const [start, end] = statusPalette[dataset.status].gradient;
      gradient.addColorStop(0, start);
      gradient.addColorStop(1, end);
      return gradient;
    },
    hoverBackgroundColor: statusPalette[dataset.status].border,
    maxBarThickness: 34
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
        color: '#64748b',
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
        color: 'rgba(148, 163, 184, 0.15)'
      },
      ticks: {
        precision: 0,
        color: '#64748b',
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
        color: '#475569',
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
    duration: 650,
    easing: 'easeOutQuart'
  }
}));
</script>

<template>
  <div class="h-80">
    <Bar :data="chartData" :options="chartOptions" />
  </div>
</template>
